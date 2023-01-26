package slack

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/xplorer-io/xplorers-bot-go/internal/ssm"
	"golang.org/x/exp/slices"
)

type SlackMessages struct {
	Messages []SlackMessage `json:"messages"`
}

type SlackMessage struct {
	Type string     `json:"type"`
	Text *SlackText `json:"text"`
}

type SlackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Slack struct {
	client *slack.Client
}

func GetSlackApiClient() (slackConfig *Slack, Error error) {
	xplorersSlackToken, err := ssm.GetSsmParameter(os.Getenv("SLACK_OAUTH_TOKEN_SSM_PATH"))
	if err != nil {
		return nil, err
	}

	return &Slack{
		client: slack.New(xplorersSlackToken),
	}, nil
}

func (s *Slack) PostSlackMessage(message string, channel string) error {
	_, _, err := s.client.PostMessage(channel, slack.MsgOptionText(message, false))
	if err != nil {
		return err
	}
	return nil
}

func (s *Slack) WelcomeNewSlackUser(slackUser string, channel string) error {
	slackMessage := GetSlackMessageBlock(slackUser)
	msgOptions := slack.MsgOptionBlocks(slackMessage...)
	if _, _, err := s.client.PostMessage(channel, msgOptions); err != nil {
		return err
	}
	return nil
}

func GetCurrentEmojiNamesOnSlackPost(GetSlackPostEmojisResponse []slack.ItemReaction) []string {
	messageReactions := []string{}
	for _, reaction := range GetSlackPostEmojisResponse {
		messageReactions = append(messageReactions, reaction.Name)
	}
	return messageReactions
}

func (s *Slack) GetCurrentEmojisOnSlackPost(msgRef slack.ItemRef) ([]slack.ItemReaction, error) {
	messageReactions, err := s.client.GetReactions(msgRef, slack.NewGetReactionsParameters())
	if err != nil {
		return nil, err
	}
	return messageReactions, nil
}

func (s *Slack) ReactToSlackPost(text string, timestamp string, channel string) error {

	emojisToReactWith := GetEmojisToReactWith(text)

	// Grab a reference to the message.
	msgRef := slack.NewRefToMessage(channel, timestamp)

	// Get current reactions on the slack post
	CurrentEmojisOnSlackPost, err := s.GetCurrentEmojisOnSlackPost(msgRef)
	if err != nil {
		return err
	}
	emojiNames := GetCurrentEmojiNamesOnSlackPost(CurrentEmojisOnSlackPost)

	// Unique slice of emojis that have not been added to the slack post
	finalEmojisToReactWith := []string{}

	for _, emoji := range emojisToReactWith {
		if !slices.Contains(emojiNames, emoji) {
			finalEmojisToReactWith = append(finalEmojisToReactWith, emoji)
		}
	}

	fmt.Printf("Reacting to slack post with emojis: %s \n", finalEmojisToReactWith)
	for _, emoji := range finalEmojisToReactWith {
		if err := s.client.AddReaction(emoji, slack.NewRefToMessage(channel, timestamp)); err != nil {
			if err.Error() == "already_reacted" {
				fmt.Printf("Already reacted to the post with emoji %s", emoji)
			} else {
				return err
			}
		}
	}
	return nil
}

func (s *Slack) ProcessSlackCallbackEvent(innerEvent slackevents.EventsAPIInnerEvent) error {
	switch ev := innerEvent.Data.(type) {

	case *slackevents.MessageEvent:
		text := ev.Text
		timestamp := ev.EventTimeStamp
		channel := ev.Channel

		switch ev.SubType {
		// Slack posts usually get edited
		case "message_changed":
			text = ev.Message.Text
			timestamp = ev.Message.TimeStamp

		case "channel_join":
			slackUser := ev.User
			fmt.Printf("%s just joined channel %s", slackUser, channel)
			if err := s.WelcomeNewSlackUser(slackUser, channel); err != nil {
				fmt.Printf("Unable to welcome new user %s with error: %s", slackUser, err)
				return err
			}
		}

		if err := s.ReactToSlackPost(text, timestamp, channel); err != nil {
			fmt.Printf("Unable to add emoji reactions to the slack post %s", err)
			return err
		}

	case *slackevents.AppMentionEvent:
		helloMessage := fmt.Sprintf("Hello <@%s> :wave: what can I do for you today?", ev.User)
		if err := s.PostSlackMessage(helloMessage, ev.Channel); err != nil {
			fmt.Printf("Unable to respond to user %s who mentioned me: %s ", ev.User, err)
			return err
		}
	}
	return nil
}
