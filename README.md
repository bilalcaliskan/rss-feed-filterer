# RSS Feed Filterer for Github Atom Feeds :newspaper_roll: :mag: :cloud:
![CI](https://github.com/bilalcaliskan/rss-feed-filterer/workflows/CI/badge.svg?event=push)
![Go Report Card](https://goreportcard.com/badge/github.com/bilalcaliskan/rss-feed-filterer)
![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_rss-feed-filterer&metric=alert_status)
![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_rss-feed-filterer&metric=sqale_rating)
![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_rss-feed-filterer&metric=reliability_rating)
![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_rss-feed-filterer&metric=security_rating)
![Coverage](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_rss-feed-filterer&metric=coverage)
![Release](https://img.shields.io/github/release/bilalcaliskan/rss-feed-filterer.svg)
![Go version](https://img.shields.io/github/go-mod/go-version/bilalcaliskan/rss-feed-filterer)
![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)

RSS Feed Filterer is a sophisticated tool designed to efficiently monitor, filter, and notify users about new releases in software projects based on their RSS feeds. It seamlessly integrates with AWS S3 to persist release data and provides a comprehensive mechanism to track multiple project releases.

## Features

- **RSS Feed Monitoring**: Efficiently checks and parses RSS feeds for software releases.
- **Notifications**: Notifies users about the latest releases. Only supports Slack notification but the architecture is designed to easily accommodate other notification services.
- **Cloud Integration**: While AWS S3 is natively supported for persistent release data storage, the architecture is designed to easily accommodate other cloud providers' S3 services in upcoming releases. (Aliyun, GCP)

## Configuration
```shell
```

## Installation
### Kubernetes
You can use [sample deployment file](resources/sample_deployment.yaml) to deploy your Kubernetes cluster.

### Binary
Binary can be downloaded from [Releases](https://github.com/bilalcaliskan/rss-feed-filterer/releases) page.

### Homebrew
This project can also be installed with [Homebrew](https://brew.sh/):
```shell
$ brew tap bilalcaliskan/tap
$ brew install bilalcaliskan/tap/rss-feed-filterer
```

## Development
This project requires below tools while developing:
- [Golang 1.20](https://golang.org/doc/go1.20)
- [pre-commit](https://pre-commit.com/)
- [golangci-lint](https://golangci-lint.run/usage/install/) - required by [pre-commit](https://pre-commit.com/)
- [gocyclo](https://github.com/fzipp/gocyclo) - required by [pre-commit](https://pre-commit.com/)

After you installed [pre-commit](https://pre-commit.com/) and the rest, simply run below command to prepare your
development environment:
```shell
$ make pre-commit-setup
```
