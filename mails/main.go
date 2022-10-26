package mails

type SMTPConfig struct {
	SMTP SMTPConfiguration
}
type Configuration struct {
	Default string
	From From
	Mailers SMTPConfig
}

// NewDefaultProvider create default provider
func NewDefaultProvider(conf Configuration) MailProvider {
	switch conf.Default {
	case "smtp":
		return NewSMTPMailer(conf.Mailers.SMTP)
	default:
		panic("Default Mail Provider is not found")
	}
}

func NewDefaultSender(conf Configuration) From {
	return conf.From
}
