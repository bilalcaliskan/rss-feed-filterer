#!/bin/bash

go test -tags "unit e2e integration" ./... -race -p 1 -coverprofile=coverage.txt -covermode=atomic
#go test -tags "unit e2e integration" ./internal/announce/email/ses -race -p 1 -coverprofile=coverage.txt -covermode=atomic
go tool cover -html=coverage.txt -o cover.html
#open cover.html
