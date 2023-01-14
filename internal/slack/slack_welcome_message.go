package slack

import (
	"math/rand"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

func GetSlackMessages() SlackMessages {
	return SlackMessages{
		Messages: []SlackMessage{
			{
				Type: "section",
				Text: &SlackText{
					Type: "mrkdwn",
					Text: "Woot Woot! We have a new *Xplorer* joining us today, welcome <@user_id> :wave: \n\n Very excited to have you join the family :charmander:",
				},
			},
			{
				Type: "section",
				Text: &SlackText{
					Type: "mrkdwn",
					Text: "Another day, another *Xplorer!* Sharing your knowledge is the most beautiful thing you guys do! \n\n Lets welcome <@user_id> to our family and lets get learning! :tada:",
				},
			},
			{
				Type: "section",
				Text: &SlackText{
					Type: "mrkdwn",
					Text: "Hey <@user_id>, welcome to the world of *Xplorers* :launch: Very excited to have you join us :smiley: \n\n We are called *Xplorers* for a reason, `We learn, share and grow together!`. \n\n Please welcome <@user_id> to the family everyone :wave: ",
				},
			},
		},
	}
}

func GetRandomMessage() SlackMessage {
	messages := GetSlackMessages().Messages
	rand.Seed(time.Now().UnixNano())
	return messages[rand.Intn(len(messages))]
}

func AddUserIdToSlackMessage(userId string, SlackMessage SlackMessage) SlackMessage {
	SlackMessage.Text.Text = strings.Replace(SlackMessage.Text.Text, "user_id", userId, -1)
	return SlackMessage
}

func GetSlackMessage(userId string) SlackMessage {
	randomMessage := GetRandomMessage()
	SlackMessage := AddUserIdToSlackMessage(userId, randomMessage)
	return SlackMessage
}

func GetSlackMessageBlock(userId string) (SlackMessageBlock []slack.Block) {
	SlackMessage := GetSlackMessage(userId)

	messageTextBlock := slack.NewTextBlockObject(slack.MarkdownType, SlackMessage.Text.Text, false, false)
	messageSectionBlock := slack.NewSectionBlock(messageTextBlock, nil, nil)

	SlackMessageBlock = append(SlackMessageBlock, messageSectionBlock)
	return SlackMessageBlock
}
