package feed

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/types"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog"
)

type Parser interface {
	ParseURL(url string) (*gofeed.Feed, error)
}

type ReleaseChecker struct {
	aws.S3ClientAPI
	bucketName string
	Parser
	logger zerolog.Logger
	config.Repository
	announcers []announce.Announcer
}

func NewReleaseChecker(client aws.S3ClientAPI, repo config.Repository, parser Parser, bucketName string, logger zerolog.Logger, announcers []announce.Announcer) *ReleaseChecker {
	return &ReleaseChecker{
		S3ClientAPI: client,
		bucketName:  bucketName,
		Parser:      parser,
		logger:      logger,
		Repository:  repo,
		announcers:  announcers,
	}
}

func (r *ReleaseChecker) CheckGithubReleases(ctx context.Context, projectName string, oneShot bool) {
	r.logger = r.logger.With().Str("projectName", projectName).Logger()

	// Run immediately since ticker does not run on first hit
	r.checkFeed(projectName, r.Repository)

	if oneShot {
		return
	}

	ticker := time.NewTicker(time.Duration(r.CheckIntervalMinutes) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.checkFeed(projectName, r.Repository)
		}
	}
}

// this method is for testing purposes
/*func (r *ReleaseChecker) fetchFeed(projectName string) (*gofeed.Feed, error) {
	r.logger.Info().Str("projectName", projectName).Msg("trying to fetch the feed")

	content, err := ioutil.ReadFile("test/releases.atom")
	if err != nil {
		log.Fatal(err)
	}
	text := string(content)
	return gofeed.NewParser().ParseString(text)
}*/

// this method is for production
func (r *ReleaseChecker) fetchFeed(projectName string) (*gofeed.Feed, error) {
	r.logger.Info().Str("projectName", projectName).Msg("trying to fetch the feed")

	// this block is for production that gets the feed from url defined in config file
	return r.ParseURL(fmt.Sprintf("%s/releases.atom", r.Url))
}

func (r *ReleaseChecker) checkFeed(projectName string, repo config.Repository) {
	for retries := 0; retries < maxRetries; retries++ {
		feed, err := r.fetchFeed(projectName)
		if err != nil {
			r.logger.Warn().Err(err).Str("url", repo.Url).Msg("an error occurred while fetching feed, retrying...")
			time.Sleep(time.Second * 5)
			continue
		}

		r.logger.Info().
			Str("name", repo.Name).
			Int("count", len(feed.Items)).
			Msgf("fetched releases")

		fetchedReleases := r.getReleasesFromFeed(projectName, feed.Items)

		var allReleases []types.Release
		if aws.IsObjectExists(r.S3ClientAPI, r.bucketName, fmt.Sprintf("%s/%s", projectName, releaseFileKey)) {
			previousReleases, err := aws.GetReleases(r.S3ClientAPI, r.bucketName, fmt.Sprintf("%s/%s", projectName, releaseFileKey))
			if err != nil {
				r.logger.Warn().Err(err).Msg("an error occured while getting releases from bucket")
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
		if err := aws.PutReleases(r.S3ClientAPI, r.bucketName, fmt.Sprintf("%s/%s", projectName, releaseFileKey), allReleases); err != nil {
			r.logger.Warn().Err(err).Msg("an error occured while putting releases into bucket")
			continue
		}

		r.logger.Info().Int("count", len(allReleases)).Msg("successfully put all releases into bucket")

		break
	}
}

func (r *ReleaseChecker) sendNotification(releases []types.Release) {
	if len(r.announcers) == 0 {
		return
	}

	for _, v := range releases {
		for _, a := range r.announcers {
			if err := a.Notify(&announce.AnnouncerPayload{
				ProjectName: v.ProjectName,
				Version:     v.Version,
				URL:         v.Url,
			}); err != nil {
				r.logger.Warn().Err(err).Msg("an error occurred while sending announce, skipping")
				continue
			}

			r.logger.Info().Str("version", v.Version).Msg("successfully sent announce")
		}
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

func (r *ReleaseChecker) getDiff(fetchedReleases []types.Release, previousReleases []types.Release) (diff []types.Release) {
	for _, item := range fetchedReleases {
		if !r.contains(previousReleases, item) {
			diff = append(diff, item)
		}
	}

	return diff
}

func (r *ReleaseChecker) extractProjectName() (string, error) {
	u, err := url.Parse(r.Url)
	if err != nil {
		return "", err
	}

	// TODO: what about other cases?
	parts := strings.Split(u.Path, "/")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid github url format")
	}

	projectName := fmt.Sprintf("%s/%s", parts[1], parts[2])
	if strings.HasPrefix(projectName, "/") || strings.HasSuffix(projectName, "/") {
		return "", fmt.Errorf("invalid project name")
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
