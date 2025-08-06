package mails

type Configuration struct {
	Default string             `json:"default"`
	SMTP    *SMTPConfiguration `json:"smtp"`
	SES     *SESConfiguration  `json:"ses"`
	From    From               `json:"from"`
}

// NewDefaultProvider create default provider
func NewDefaultProvider(conf Configuration) MailProvider {
	switch conf.Default {
	case "smtp":
		return NewSMTPMailer(*conf.SMTP, conf.From)
	case "ses":
		return NewSESProvider(*conf.SES, conf.From)
	default:
		panic("Default Mail Provider is not found")
	}
}
