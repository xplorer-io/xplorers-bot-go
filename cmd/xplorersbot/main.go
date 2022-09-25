package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/getsentry/sentry-go"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/xplorer-io/xplorers-bot-go/internal/xplorersbot"
)

func ProcessSlackCallbackEvent(innerEvent slackevents.EventsAPIInnerEvent, apiClient *slack.Client) {
	switch ev := innerEvent.Data.(type) {

	case *slackevents.MessageEvent:
		if ev.SubType == "channel_join" {
			fmt.Printf("%s just joined channel %s", ev.User, ev.Channel)
			welcomeMessage := xplorersbot.GetWelcomeMessageBlock(ev.User)
			_, _, postSlackMessageErr := apiClient.PostMessage(ev.Channel, slack.MsgOptionBlocks(welcomeMessage...))
			if postSlackMessageErr != nil {
				fmt.Printf("Unable to welcome new user %s with error: %s", ev.User, postSlackMessageErr)
				sentry.CaptureException(postSlackMessageErr)
			}
		}

		emojisToAdd := xplorersbot.GetEmojis(ev.Text)
		if len(emojisToAdd) > 0 {
			fmt.Println("Reacting to slack post with emojis: ", emojisToAdd)
			for _, emoji := range emojisToAdd {
				addReactionErr := apiClient.AddReaction(emoji, slack.NewRefToMessage(ev.Channel, ev.EventTimeStamp))
				if addReactionErr != nil {
					fmt.Printf("Unable to add emoji reaction `%s` to the slack post", emoji)
					sentry.CaptureException(addReactionErr)
				}
			}
		}

	case *slackevents.AppMentionEvent:
		helloMessage := fmt.Sprintf("Hello <@%s> :wave: what can I do for you today?", ev.User)
		_, _, postSlackMessageErr := apiClient.PostMessage(ev.Channel, slack.MsgOptionText(helloMessage, false))
		if postSlackMessageErr != nil {
			fmt.Printf("Unable to respond to user %s who mentioned me: %s ", ev.User, postSlackMessageErr)
			sentry.CaptureException(postSlackMessageErr)
		}
	}
}

func ProcessSlackEvent(request events.APIGatewayProxyRequest, eventsAPIEvent slackevents.EventsAPIEvent, apiClient *slack.Client) (events.APIGatewayProxyResponse, error) {
	switch eventsAPIEvent.Type {

	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		jsonUnmarshalErr := json.Unmarshal([]byte(request.Body), &r)
		if jsonUnmarshalErr != nil {
			fmt.Println("Unable to unmarshal slack url verification event.")
			sentry.CaptureException(jsonUnmarshalErr)
		}
		fmt.Println("URL verification event, responding with the challenge.")
		return events.APIGatewayProxyResponse{Body: r.Challenge, StatusCode: 200}, nil

	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		ProcessSlackCallbackEvent(innerEvent, apiClient)
	}

	// Ignore any events that dont match our rules
	return events.APIGatewayProxyResponse{Body: "Successfully processed slack event", StatusCode: 200}, nil
}

func ProcessSlackRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sentryDsnParameterPath := os.Getenv("SENTRY_DSN_SSM_PATH")
	sentryDsn, getParamErr := xplorersbot.GetSsmParameter(&sentryDsnParameterPath)
	if getParamErr != nil {
		fmt.Printf("Unable to fetch sentry dsn: %s", getParamErr)
	}

	sentryErr := sentry.Init(sentry.ClientOptions{
		Dsn: sentryDsn,

		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,

		// Either set environment and release here or set the SENTRY_ENVIRONMENT
		// and SENTRY_RELEASE environment variables.
		Environment: os.Getenv("ENVIRONMENT"),
		// Release:     "my-project-name@1.0.0",

		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,
	})
	if sentryErr != nil {
		log.Fatalf("sentry.InitError: %s", sentryErr)
	}
	defer sentry.Flush(time.Second)

	fmt.Printf("Slack event body: %s\n", request.Body)

	apiClient, getApiClientErr := xplorersbot.GetSlackApiClient()
	if getApiClientErr != nil {
		fmt.Println("Unable to instantiate slack api client")
		sentry.CaptureException(getApiClientErr)
	}

	eventsAPIEvent, parseSlackEventErr := slackevents.ParseEvent(json.RawMessage(request.Body), slackevents.OptionNoVerifyToken())

	if parseSlackEventErr != nil {
		fmt.Println("Failed to parse slack event body.")
		sentry.CaptureException(parseSlackEventErr)
		return events.APIGatewayProxyResponse{Body: "Bad request body!", StatusCode: 400}, nil
	}

	return ProcessSlackEvent(request, eventsAPIEvent, apiClient)
}

func main() {
	lambda.Start(ProcessSlackRequest)
}
