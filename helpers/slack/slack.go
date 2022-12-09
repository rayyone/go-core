package plugins

import (
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"sync"
)

type SlackConfig struct {
	Token string
	DefaultChannel string
	Environment string
}
type Slack struct {
	api *slack.Client
	config SlackConfig
	options map[string]interface{}
	optionRWLock sync.RWMutex
}
var (
	currentSlackClient = NewSlackClient("", SlackConfig{})
)

func CurrentSlackClient() *Slack {
	return currentSlackClient
}
func NewSlackClient(token string, config SlackConfig) *Slack {
	api := slack.New(token)
	return &Slack{
		api: api,
		config: config,
		options: map[string]interface{}{},
		optionRWLock: sync.RWMutex{},
	}
}

func InitSlackClient(config SlackConfig) error {
	slackClient := CurrentSlackClient()
	slackClient.api = slack.New(config.Token)
	slackClient.config = config
	return nil
}
func (s *Slack) SetOption(key string, value interface{}) {
	s.optionRWLock.RLock()
	s.options[key] = value
	s.optionRWLock.RUnlock()
}

func (s *Slack) GetOption(key string) interface{} {
	val, ok := s.options[key]
	if !ok {
		return nil
	}
	return val
}

func (s *Slack) SendSimpleMessageToChannel(channel string, title string, message string) error {
	if s.config.Token == "" {
		return errors.New("slack client not init")
	}
	if s.config.Environment != "production" {
		title += fmt.Sprintf(" (%s) ENV", s.config.Environment)
	}
	channelID, timestamp, err := s.api.PostMessage(
		channel,
		slack.MsgOptionText(fmt.Sprintf("*%s*\n\n%s", title, message), false),
		slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	return nil
}
func SendSimpleMessage(title string, message string) error {
	slackClient := CurrentSlackClient()
	if slackClient.config.Token == "" {
		fmt.Printf("slack client not init")
		return nil
	}
	return slackClient.SendSimpleMessageToChannel(slackClient.config.DefaultChannel, title, message)
}
func SendSimpleMessageToChannel(channel string, title string, message string) error {
	slackClient := CurrentSlackClient()
	if slackClient.config.Token == "" {
		fmt.Printf("slack client not init")
		return nil
	}
	return slackClient.SendSimpleMessageToChannel(channel, title, message)
}
