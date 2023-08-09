package feed

import (
	"context"
	"fmt"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/notification/slack"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/types"

	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog"
)

const (
	maxRetries         = 3
	defaultSemverRegex = `/(v?\d+\.\d+\.\d+)$`
	releaseFileKey     = "releases.json"
)

type ReleaseChecker struct {
	client     aws.S3ClientAPI
	bucketName string
	parser     *gofeed.Parser
	logger     zerolog.Logger
	repo       config.Repository
}

func NewReleaseChecker(client aws.S3ClientAPI, repo config.Repository, bucketName string, logger zerolog.Logger) *ReleaseChecker {
	return &ReleaseChecker{
		client:     client,
		bucketName: bucketName,
		parser:     gofeed.NewParser(),
		logger:     logger,
		repo:       repo,
	}
}

func (r *ReleaseChecker) CheckGithubReleases(ctx context.Context, sem chan struct{}) {
	projectName, err := r.extractProjectName(r.repo.Url)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to extract project name")
		return
	}

	r.logger = r.logger.With().Str("projectName", projectName).Logger()

	ticker := time.NewTicker(time.Duration(r.repo.CheckIntervalMinutes) * time.Minute)
	defer ticker.Stop()

	// Run immediately since ticker does not run on first hit
	r.checkFeed(sem, projectName, r.repo)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.checkFeed(sem, projectName, r.repo)
		}
	}
}

func (r *ReleaseChecker) checkFeed(sem chan struct{}, projectName string, repo config.Repository) {
	defer func() {
		<-sem // release the semaphore
	}()

	for retries := 0; retries < maxRetries; retries++ {
		r.logger.Info().Str("projectName", projectName).Msg("trying to fetch the feed")

		file, err := os.Open("testdata/releases.atom")
		if err != nil {
			r.logger.Warn().Str("error", err.Error()).Msg("Error opening file")
			continue
		}

		feed, err := r.parser.Parse(file)
		if err != nil {
			r.logger.Warn().
				Str("error", err.Error()).
				Msg("an error occurred while parsing feed, retrying...")
			continue
		}

		_ = file.Close()

		//feed, err := r.parser.ParseURL(fmt.Sprintf("%s/releases.atom", repo.Url))
		//if err != nil {
		//	r.logger.Warn().
		//		Str("error", err.Error()).
		//		Str("url", repo.Url).
		//		Msg("an error occurred while fetching feed, retrying...")
		//	time.Sleep(time.Second * 5)
		//	continue
		//}

		r.logger.Info().
			Str("name", repo.Name).
			Int("count", len(feed.Items)).
			Msgf("fetched releases")

		fetchedReleases := r.getReleasesFromFeed(projectName, feed.Items)

		var allReleases []types.Release
		if aws.IsObjectExists(r.client, r.bucketName, fmt.Sprintf("%s/%s", projectName, releaseFileKey)) {
			previousReleases, err := aws.GetReleases(r.client, r.bucketName, fmt.Sprintf("%s/%s", projectName, releaseFileKey))
			if err != nil {
				r.logger.Warn().Msg("an error occured while getting releases from bucket")
				continue
			}

			diff := r.getDiff(fetchedReleases, previousReleases)
			if len(diff) == 0 {
				r.logger.Info().Msg("no new releases found, nothing to do")
				return
			}

			r.logger.Info().Int("count", len(diff)).Msg("successfully fetched diffs")
			r.sendNotification(diff)

			allReleases = append(diff, previousReleases...)
		} else {
			r.logger.Info().Msg("releases does not exists on bucket, adding from scratch")
			allReleases = fetchedReleases
		}

		r.logger.Info().Msg("putting diffs into bucket")
		if err := aws.PutReleases(r.client, r.bucketName, fmt.Sprintf("%s/%s", projectName, releaseFileKey), allReleases); err != nil {
			r.logger.Warn().Msg("an error occured while putting releases into bucket")
			continue
		}

		r.logger.Info().Int("count", len(allReleases)).Msg("successfully put all releases into bucket")

		// TODO: check if notifications enabled
		// TODO: ensure project name does not end with /

		//r.sendNotification(diff)

		break
	}
}

func (r *ReleaseChecker) sendNotification(releases []types.Release) {
	for _, v := range releases {
		if err := slack.SendNotification(v.ProjectName, v.Version, v.Url); err != nil {
			r.logger.Warn().Err(err).Msg("an error occurred while sending notification, skipping")
			continue
		}

		r.logger.Info().Msg("successfully sent notification")
	}
}

func (r *ReleaseChecker) getReleasesFromFeed(projectName string, items []*gofeed.Item) []types.Release {
	var releases []types.Release
	for _, item := range items {
		pattern := regexp.MustCompile(defaultSemverRegex)
		matches := pattern.FindStringSubmatch(item.Link)
		if len(matches) > 0 {
			releases = append(releases, types.Release{
				ProjectName: projectName,
				Version:     item.Title,
				PublishedAt: item.PublishedParsed,
				UpdatedAt:   item.UpdatedParsed,
				Url:         item.Link,
			})
		}
	}

	return releases
}

//func printReleases(releases []types.Release) {
//	for _, v := range releases {
//		fmt.Printf("%s %s\n", v.ProjectName, v.Version)
//	}
//}

func (r *ReleaseChecker) getDiff(fetchedReleases []types.Release, previousReleases []types.Release) (diff []types.Release) {
	for _, item := range fetchedReleases {
		if !r.contains(previousReleases, item) {
			diff = append(diff, item)
		}
	}

	return diff
}

func (r *ReleaseChecker) extractProjectName(repoUrl string) (string, error) {
	u, err := url.Parse(repoUrl)
	if err != nil {
		return "", err
	}

	parts := strings.Split(u.Path, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid github url format")
	}

	return fmt.Sprintf("%s/%s", parts[1], parts[2]), nil
}

func (r *ReleaseChecker) contains(releases []types.Release, release types.Release) bool {
	for _, item := range releases {
		if release.ProjectName == item.ProjectName &&
			release.Version == item.Version &&
			release.Url == item.Url &&
			(release.PublishedAt == nil && item.PublishedAt == nil ||
				release.PublishedAt != nil && item.PublishedAt != nil && release.PublishedAt.Equal(*item.PublishedAt)) &&
			(release.UpdatedAt == nil && item.UpdatedAt == nil ||
				release.UpdatedAt != nil && item.UpdatedAt != nil && release.UpdatedAt.Equal(*item.UpdatedAt)) {
			return true
		}
	}

	return false
}
