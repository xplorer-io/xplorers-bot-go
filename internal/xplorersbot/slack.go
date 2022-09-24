package xplorersbot

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func GetSlackApiClient() (slackApiClient *slack.Client, Error error) {
	xplorersBotTokenSsmParameterPath := os.Getenv("SLACK_OAUTH_TOKEN_SSM_PATH")
	xplorersBotToken, err := GetSsmParameter(&xplorersBotTokenSsmParameterPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return slack.New(xplorersBotToken), nil
}
