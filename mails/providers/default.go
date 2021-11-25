package providers

import (
	"github.com/rayyone/go-core/mails"
	"github.com/rayyone/go-core/mails/providers/smtp"
)

type Configuration struct {
	Default string
	SMTP    smtp.Configuration
}

// NewDefaultProvider create default provider
func NewDefaultProvider(conf Configuration) mails.MailProvider {
	switch conf.Default {
	case "smtp":
		return smtp.NewMailer(conf.SMTP)
	default:
		panic("Default Mail Provider is not found")
	}
}
