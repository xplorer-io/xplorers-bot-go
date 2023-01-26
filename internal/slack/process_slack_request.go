package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/getsentry/sentry-go"
	"github.com/slack-go/slack/slackevents"
	"github.com/xplorer-io/xplorers-bot-go/internal/sentrylib"
	"github.com/xplorer-io/xplorers-bot-go/internal/xplorersboterrors"
)

func ProcessRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if err := sentrylib.InitializeSentry(); err != nil {
		log.Fatalf("could not initialise sentry: %v", err)
	}
	defer sentry.Flush(time.Second)
	fmt.Printf("Slack event body: %s\n", request.Body)

	slackConfig, err := GetSlackApiClient()
	if err != nil {
		fmt.Printf("Unable to instantiate slack %s", err)
		sentry.CaptureException(err)
		return events.APIGatewayProxyResponse{Body: xplorersboterrors.InternalServerError.String(), StatusCode: http.StatusOK}, nil
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(request.Body), slackevents.OptionNoVerifyToken())
	if err != nil {
		fmt.Printf("Failed to parse slack event body %s", err)
		sentry.CaptureException(err)
		return events.APIGatewayProxyResponse{Body: xplorersboterrors.InternalServerError.String(), StatusCode: http.StatusOK}, nil
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		if err := json.Unmarshal([]byte(request.Body), &r); err != nil {
			fmt.Println("Unable to unmarshal slack url verification event.")
			break
		}
		fmt.Println("URL verification event, responding with the challenge.")
		return events.APIGatewayProxyResponse{Body: r.Challenge, StatusCode: 200}, nil

	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		if err := slackConfig.ProcessSlackCallbackEvent(innerEvent); err != nil {
			fmt.Printf("Unable to process callback event: %s ", err)
		}
	}

	return events.APIGatewayProxyResponse{Body: "Successfully processed slack event", StatusCode: http.StatusOK}, nil
}
