package plugins

import (
	"fmt"
	"github.com/slack-go/slack"
)
type SlackConfig struct {
	Token string
	DefaultChannel string
	Environment string
}
type Slack struct {
	api *slack.Client
	config SlackConfig
}

var currentSlackClient = NewSlackClient("", SlackConfig{})
func CurrentSlackClient() *Slack {
	return currentSlackClient
}
func NewSlackClient(token string, config SlackConfig) *Slack {
	api := slack.New(token)
	return &Slack{
		api: api,
		config: config,
	}
}

func InitSlackClient(config SlackConfig) error {
	slackClient := CurrentSlackClient()
	slackClient.api = slack.New(config.Token)
	slackClient.config = config
	return nil
}

func SendSimpleMessage(title string, message string) {
	slackClient := CurrentSlackClient()
	SendSimpleMessageToChannel(slackClient.config.DefaultChannel, title, message)
}

func SendSimpleMessageToChannel(channel string, title string, message string) {
	slackClient := CurrentSlackClient()
	if slackClient.config.Token == "" {
		return
	}
	if slackClient.config.Environment != "production" {
		title += fmt.Sprintf(" (%s) ENV", slackClient.config.Environment)
	}
	channelID, timestamp, err := slackClient.api.PostMessage(
		channel,
		slack.MsgOptionText(fmt.Sprintf("*%s*\n\n%s", title, message), false),
		slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
}
