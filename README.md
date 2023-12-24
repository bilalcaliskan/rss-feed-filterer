# RSS Feed Filterer for Github Atom Feeds :newspaper_roll: :mag: :cloud:
[![CI](https://github.com/bilalcaliskan/rss-feed-filterer/actions/workflows/push.yml/badge.svg)](https://github.com/bilalcaliskan/rss-feed-filterer/actions/workflows/push.yml)
![Go Report Card](https://goreportcard.com/badge/github.com/bilalcaliskan/rss-feed-filterer)
![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_rss-feed-filterer&metric=alert_status)
![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_rss-feed-filterer&metric=sqale_rating)
![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_rss-feed-filterer&metric=reliability_rating)
![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_rss-feed-filterer&metric=security_rating)
![Coverage](https://sonarcloud.io/api/project_badges/measure?project=bilalcaliskan_rss-feed-filterer&metric=coverage)
![Release](https://img.shields.io/github/release/bilalcaliskan/rss-feed-filterer.svg)
![Go version](https://img.shields.io/github/go-mod/go-version/bilalcaliskan/rss-feed-filterer)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit)](https://github.com/pre-commit/pre-commit)
![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)

RSS Feed Filterer is a sophisticated tool designed to efficiently monitor, filter, and notify users about new releases in software projects based on their RSS feeds. It seamlessly integrates with AWS S3 to persist release data and provides a comprehensive mechanism to track multiple project releases.

## Features

- **RSS Feed Monitoring**: Efficiently checks and parses RSS feeds for software releases.
- **Notifications**: Notifies users about the latest releases. Only supports Slack notification but the architecture is designed to easily accommodate other notification services like email.
- **Cloud Integration**: While **AWS S3** is natively supported for persistent release data storage, the architecture is designed to easily accommodate other cloud providers' S3 services in upcoming releases. (Aliyun, GCP)

## Configuration
```shell
Usage:
  rss-feed-filterer [flags]
  rss-feed-filterer [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       starts the main process by reading the config file

Flags:
  -c, --config-file string   path for the config file to be used
  -h, --help                 help for rss-feed-filterer
      --verbose              verbose output of the logging library as 'debug' (default false)
  -v, --version              version for rss-feed-filterer

Use "rss-feed-filterer [command] --help" for more information about a command.
```

## Installation
### Kubernetes
You can use [sample deployment file](deployments/sample_valid_deployment.yaml) to deploy your Kubernetes cluster.

### Binary
Binary can be downloaded from [Releases](https://github.com/bilalcaliskan/rss-feed-filterer/releases) page.

## Development
This project requires below tools while developing:
- [Golang 1.21](https://golang.org/doc/go1.21)
- [pre-commit](https://pre-commit.com/)

After you installed [pre-commit](https://pre-commit.com/) and the rest, simply run below command to prepare your
development environment:
```shell
$ make pre-commit-setup
```
