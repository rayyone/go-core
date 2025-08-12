package mails

import (
	"bytes"
	"fmt"
	"github.com/go-mail/mail"
	"github.com/rayyone/go-core/ryerr"
)

type Provider struct{}

func (*Provider) BuildMailMessage(msg Message) (*mail.Message, error) {
	if msg.From == nil || msg.From.Name == "" {
		return nil, ryerr.New("From address is empty")
	}
	fromStr := msg.From.Address
	if msg.From.Name != "" {
		fromStr = fmt.Sprintf("%s <%s>", msg.From.Name, msg.From.Address)
	}
	mailMsg := mail.NewMessage()
	mailMsg.SetHeader("From", fromStr)
	mailMsg.SetHeader("To", msg.To...)
	mailMsg.SetHeader("Cc", msg.Cc...)
	mailMsg.SetHeader("Bcc", msg.Bcc...)
	mailMsg.SetHeader("Subject", msg.Subject)
	if msg.Header != nil {
		for k, v := range msg.Header {
			mailMsg.SetHeader(k, v)
		}
	}
	if msg.Text != "" {
		mailMsg.SetBody("text/plain", msg.Text)
	}
	if msg.HTML != "" {
		mailMsg.AddAlternative("text/html", msg.HTML)
	}

	for _, attachment := range msg.Attachments {
		mailMsg.AttachReader(
			attachment.Name,
			bytes.NewReader(attachment.Content),
			mail.SetHeader(map[string][]string{
				"Content-Type": {attachment.ContentType},
			}),
		)
	}
	return mailMsg, nil
}
