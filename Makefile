SHELL = /bin/bash
SHELLFLAGS = -ex

VERSION ?= $(shell git rev-parse --short HEAD)
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
API_STAGE ?= dev
ENVIRONMENT ?= dev

# Import settings and stage-specific overrides
include ./settings/defaults.conf
ifneq ("$(wildcard ./settings/$(ENVIRONMENT).conf"), "")
-include ./settings/$(ENVIRONMENT).conf
endif

help:  ## get help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

download-golangci-lint: ## download golang ci linter
	@echo "--- downloading golang ci linter ---"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${GOLANGCI_VERSION}

lint: download-golangci-lint ## run golang linter
	@echo "--- running golang linter ---"
	@bin/golangci-lint run

lint-fix: ## run golang linter with fix option
	@echo "--- running golang linter with fix option ---"
	@bin/golangci-lint run --fix

build: ## build go binary
	@echo "--- building binary using go build ---"
	@go mod tidy
	# the -w -s flags make the binary a bit smaller and
	# trimpath shortens build paths in stack traces
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags='-w -s' -trimpath -o dist/xplorersbot ./cmd/xplorersbot

test: ## run go tests and generate coverage
	@echo "--- running go tests ---"
	@mkdir -p coverage
	@go test -coverprofile=coverage/coverage.txt -covermode count ./...
	@go tool cover -func coverage/coverage.txt

zip-artifact: build ## zip binary artifact to be uploaded to s3
	@echo "--- zip binary artifact ---"
	@cd dist && zip -r -q ../lambda.zip .

cleanup: ## remove stale artifacts
	@echo "--- removing stale artifacts ---"
	@rm -rf dist/
	@rm -rf lambda.zip
	@rm -rf xplorersbot.packaged.yml

deploy-xplorers-bot: cleanup zip-artifact ## deploy xplorers bot to aws
	@echo "--- Packaging xplorersbot artifact to S3 ---"
	$(eval S3_BUCKET := $(shell aws ssm get-parameter --name $(ARTIFACTS_BUCKET_SSM_PATH) --query Parameter.Value --output text))
	aws cloudformation package \
		--s3-bucket $(S3_BUCKET) \
		--s3-prefix xplorersbot/$(GIT_BRANCH)/$(VERSION) \
		--template-file cloudformation/xplorersbot.yml \
		--output-template-file xplorersbot.packaged.yml
	aws cloudformation deploy \
		--s3-bucket $(S3_BUCKET) \
		--s3-prefix xplorersbot/$(GIT_BRANCH)/$(VERSION) \
		--template-file xplorersbot.packaged.yml \
		--stack-name xplorersbot-$(GIT_BRANCH)-deploy \
		--capabilities CAPABILITY_NAMED_IAM \
		--no-fail-on-empty-changeset \
		--parameter-overrides \
			StageName=$(API_STAGE) \
			SlackOauthTokenSsmPath=$(SLACK_OAUTH_TOKEN_SSM_PATH) \
			SentryDsnSsmPath=$(SENTRY_DSN_SSM_PATH) \
			Environment=$(ENVIRONMENT) \
		--tags xplorersbot:version=$(VERSION) xplorersbot:branch=$(GIT_BRANCH)
	@make cleanup

.PHONY: deploy deploy-xplorers-bot cleanup zip-artifact test lint-fix lint download-golangci-lint help
