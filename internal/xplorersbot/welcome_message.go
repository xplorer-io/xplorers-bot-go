package xplorersbot

import (
	"math/rand"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

type WelcomeMessages struct {
	Messages []WelcomeMessage `json:"messages"`
}

type WelcomeMessage struct {
	Type string       `json:"type"`
	Text *WelcomeText `json:"text"`
}

type WelcomeText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func GetWelcomeMessages() WelcomeMessages {
	return WelcomeMessages{
		Messages: []WelcomeMessage{
			{
				Type: "section",
				Text: &WelcomeText{
					Type: "mrkdwn",
					Text: "Woot Woot! We have a new *Xplorer* joining us today, welcome <@user_id> :wave: \n\n Very excited to have you join the family :charmander:",
				},
			},
			{
				Type: "section",
				Text: &WelcomeText{
					Type: "mrkdwn",
					Text: "Another day, another *Xplorer!* Sharing your knowledge is the most beautiful thing you guys do! \n\n Lets welcome <@user_id> to our family and lets get learning! :tada:",
				},
			},
			{
				Type: "section",
				Text: &WelcomeText{
					Type: "mrkdwn",
					Text: "Hey <@user_id>, welcome to the world of *Xplorers* :launch: Very excited to have you join us :smiley: \n\n We are called *Xplorers* for a reason, `We learn, share and grow together!`. \n\n Please welcome <@user_id> to the family everyone :wave: ",
				},
			},
		},
	}
}

func GetRandomMessage() WelcomeMessage {
	messages := GetWelcomeMessages().Messages
	rand.Seed(time.Now().UnixNano())
	return messages[rand.Intn(len(messages))]
}

func AddUserIdToWelcomeMessage(userId string, welcomeMessage WelcomeMessage) WelcomeMessage {
	welcomeMessage.Text.Text = strings.Replace(welcomeMessage.Text.Text, "user_id", userId, -1)
	return welcomeMessage
}

func GetWelcomeMessage(userId string) WelcomeMessage {
	randomMessage := GetRandomMessage()
	welcomeMessage := AddUserIdToWelcomeMessage(userId, randomMessage)
	return welcomeMessage
}

func GetWelcomeMessageBlock(userId string) (welcomeMessageBlock []slack.Block) {
	welcomeMessage := GetWelcomeMessage(userId)

	messageTextBlock := slack.NewTextBlockObject(slack.MarkdownType, welcomeMessage.Text.Text, false, false)
	messageSectionBlock := slack.NewSectionBlock(messageTextBlock, nil, nil)

	welcomeMessageBlock = append(welcomeMessageBlock, messageSectionBlock)
	return welcomeMessageBlock
}
