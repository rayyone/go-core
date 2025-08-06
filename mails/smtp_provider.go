package mails

import (
	"github.com/go-mail/mail"
	loghelper "github.com/rayyone/go-core/helpers/log"
	"github.com/rayyone/go-core/helpers/method"
	"github.com/rayyone/go-core/ryerr"
	"strings"
	"time"
)

type SMTPConfiguration struct {
	Host     string
	Port     int
	User     string
	Password string
	From     *From
}
type SMTPProvider struct {
	Provider
	config SMTPConfiguration
}

func NewSMTPMailer(conf SMTPConfiguration, from From) *SMTPProvider {
	if conf.From == nil || (conf.From != nil && conf.From.Name == "") {
		conf.From = method.Ptr(from)
	}
	return &SMTPProvider{config: conf}
}

func (s *SMTPProvider) Send(msg Message) error {

	port := 587
	if s.config.Port != 0 {
		port = s.config.Port
	}
	if msg.From == nil || msg.From.Address == "" {
		msg.From = s.config.From
	}
	mailMsg, err := s.BuildMailMessage(msg)
	if err != nil {
		return err
	}
	dialer := mail.Dialer{
		Host:         s.config.Host,
		Port:         port,
		Username:     s.config.User,
		Password:     s.config.Password,
		SSL:          port == 465,
		Timeout:      30 * time.Second,
		RetryFailure: true,
	}
	start := time.Now()
	loghelper.PrintYellowf(
		"[SMTP] Sending email to: %s, cc: %s, bcc: %s",
		strings.Join(msg.To, ", "),
		strings.Join(msg.Cc, ", "),
		strings.Join(msg.Bcc, ", "),
	)
	if err := dialer.DialAndSend(mailMsg); err != nil {
		loghelper.PrintRedf("[SMTP] Send email completed with error in %.2fs", time.Since(start).Seconds())
		return ryerr.New(ryerr.Wrapf(err, err.Error()).Error())
	}
	loghelper.PrintYellowf("[SMTP] Send email completed in %.2fs", time.Since(start).Seconds())

	return nil
}
