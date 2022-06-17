package mails

type Configuration struct {
	Default string
	SMTP    SMTPConfiguration
}

// NewDefaultProvider create default provider
func NewDefaultProvider(conf Configuration) MailProvider {
	switch conf.Default {
	case "smtp":
		return NewSMTPMailer(conf.SMTP)
	default:
		panic("Default Mail Provider is not found")
	}
}
