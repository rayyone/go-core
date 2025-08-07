package mails

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	loghelper "github.com/rayyone/go-core/helpers/log"
	"github.com/rayyone/go-core/helpers/method"
	"github.com/rayyone/go-core/ryerr"
	"strings"
	"time"
)

type SESConfiguration struct {
	Profile   string
	AccessKey string
	SecretKey string
	Region    string
	From      *From
}
type SESProvider struct {
	Provider
	config SESConfiguration
	client *sesv2.Client
}

func NewSESProvider(conf SESConfiguration, from From) *SESProvider {
	var optFn config.LoadOptionsFunc
	if conf.From == nil || (conf.From != nil && conf.From.Name == "") {
		conf.From = method.Ptr(from)
	}
	if conf.Profile != "" {
		optFn = config.WithSharedConfigProfile(conf.Profile)
	}
	if conf.AccessKey != "" && conf.SecretKey != "" {
		creds := credentials.NewStaticCredentialsProvider(conf.AccessKey, conf.SecretKey, "")
		optFn = config.WithCredentialsProvider(creds)
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(conf.Region),
		optFn,
	)
	if err != nil {
		panic(err)
	}

	client := sesv2.NewFromConfig(cfg)
	return &SESProvider{client: client, config: conf}
}

func (s *SESProvider) Send(msg Message) error {
	if msg.From == nil || msg.From.Address == "" {
		msg.From = s.config.From
	}
	mailMsg, err := s.BuildMailMessage(msg)
	if err != nil {
		return err
	}
	var buf bytes.Buffer

	_, err = mailMsg.WriteTo(&buf)
	if err != nil {
		return err
	}
	start := time.Now()
	loghelper.PrintYellowf(
		"[SMTP] Sending email to: %s, cc: %s, bcc: %s",
		strings.Join(msg.To, ", "),
		strings.Join(msg.Cc, ", "),
		strings.Join(msg.Bcc, ", "),
	)
	input := &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Raw: &types.RawMessage{
				Data: buf.Bytes(),
			},
		},
	}
	_, err = s.client.SendEmail(context.Background(), input)
	if err != nil {
		loghelper.PrintRedf("[SMTP] Send email completed with error in %.2fs", time.Since(start).Seconds())
		return ryerr.New(ryerr.Wrapf(err, err.Error()).Error())
	}
	loghelper.PrintYellowf("[SMTP] Send email completed in %.2fs", time.Since(start).Seconds())
	return err
}
