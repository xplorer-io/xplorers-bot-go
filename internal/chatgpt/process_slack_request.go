package chatgpt

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/getsentry/sentry-go"
	"github.com/slack-go/slack"
	"github.com/xplorer-io/xplorers-bot-go/internal/common"
	"github.com/xplorer-io/xplorers-bot-go/internal/sentrylib"
	"github.com/xplorer-io/xplorers-bot-go/internal/ssm"
	"github.com/xplorer-io/xplorers-bot-go/internal/xplorersboterrors"
)

func ProcessSlackRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if err := sentrylib.InitializeSentry(); err != nil {
		log.Fatalf("could not initialise sentry: %v", err)
	}
	defer sentry.Flush(time.Second)

	slackEvent, err := common.ParseQueryString(request.Body)
	fmt.Printf("Slack event body %s", slackEvent)
	if err != nil {
		fmt.Printf("Failed to parse slack event body %s", err)
		sentry.CaptureException(err)
		return events.APIGatewayProxyResponse{Body: xplorersboterrors.InternalServerError.String(), StatusCode: http.StatusOK}, nil
	}

	xplorersChatbotSlackChannelId, err := ssm.GetSsmParameter(os.Getenv("XPLORERS_CHATBOT_SLACK_CHANNEL_ID_SSM_PATH"))
	if err != nil {
		fmt.Printf("Failed to get xplorers chatbot slack channel id %s", err)
		sentry.CaptureException(err)
		return events.APIGatewayProxyResponse{Body: xplorersboterrors.InternalServerError.String(), StatusCode: http.StatusOK}, nil
	}

	slackChannelId := slackEvent.Get("channel_id")
	if slackChannelId == xplorersChatbotSlackChannelId {
		chatGPT, err := NewChatGPTAPI()
		if err != nil {
			fmt.Printf("Unable to initialize chatgpt: %s", err)
			sentry.CaptureException(err)
			return events.APIGatewayProxyResponse{Body: xplorersboterrors.InternalServerError.String(), StatusCode: http.StatusOK}, nil
		}
		slackResponseBlocks, err := chatGPT.AskXplorersBot(slackEvent)
		if err != nil {
			fmt.Printf("Unable to get a response from XplorersBot: %s", err)
			sentry.CaptureException(err)
			return events.APIGatewayProxyResponse{Body: xplorersboterrors.InternalServerError.String(), StatusCode: http.StatusOK}, nil
		}

		if err := slack.PostWebhook(slackEvent.Get("response_url"), &slack.WebhookMessage{Blocks: slackResponseBlocks, ResponseType: "in_channel"}); err != nil {
			fmt.Printf("Unable to respond to user with chatgpt response: %s", err)
			sentry.CaptureException(err)
		}

	} else {
		fmt.Printf("Slack user %s invoked slash command /chatbot from an unauthorized slack channel %s\n", slackEvent.Get("user_name"), slackChannelId)
		sentry.CaptureException(err)
	}

	return events.APIGatewayProxyResponse{Body: "Successfully processed slack event", StatusCode: http.StatusOK}, nil
}
