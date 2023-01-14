package chatgpt

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/slack-go/slack"
	islack "github.com/xplorer-io/xplorers-bot-go/internal/slack"
	"github.com/xplorer-io/xplorers-bot-go/internal/ssm"
)

type ChatGPT struct {
	client gpt3.Client
}

func NewChatGPTAPI() (*ChatGPT, error) {
	chatgptApiKey, err := ssm.GetSsmParameter(os.Getenv("CHATGPT_API_KEY_SSM_PATH"))
	if err != nil {
		return nil, fmt.Errorf("unable to fetch chat GPT API key: %v", err)
	}

	return &ChatGPT{
		client: gpt3.NewClient(chatgptApiKey),
	}, nil
}

func (c *ChatGPT) AskXplorersBot(slackEvent url.Values) (*slack.Blocks, error) {
	chatGptResponse, err := c.AskChatGPT(slackEvent.Get("text"))
	if err != nil {
		return nil, fmt.Errorf("unable to interact with chat GPT: %v", err)
	}

	slackResponse := GetFormattedChatGptResponseForSlack(chatGptResponse)
	fmt.Printf("Response from ChatGpt %s \n", chatGptResponse)

	return &slack.Blocks{BlockSet: []slack.Block{slackResponse}}, nil
}

func (c *ChatGPT) AskChatGPT(prompt string) (string, error) {
	ctx := context.Background()

	resp, err := c.client.CompletionWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt:    []string{prompt},
		MaxTokens: gpt3.IntPtr(1000),
		Echo:      true,
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Text, nil
}

func GetFormattedChatGptResponseForSlack(responseText string) *slack.SectionBlock {
	chatGptResponseMessage := islack.SlackMessage{
		Type: "section",
		Text: &islack.SlackText{
			Type: "mrkdwn",
			Text: responseText,
		},
	}

	messageTextBlock := slack.NewTextBlockObject(slack.MarkdownType, chatGptResponseMessage.Text.Text, false, false)
	return slack.NewSectionBlock(messageTextBlock, nil, nil)
}
