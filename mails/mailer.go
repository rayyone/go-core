package mails

// Mailable represent mail interface
type Mailable interface {
	BuildSubject() string
	BuildHTMLBody() string
	BuildTextBody() string
}

type MailProvider interface {
	Send(recipient Recipient, from From, content MailContent) error
	SendWithCalendarEvent(recipient Recipient, from From, content MailContent, options *CalendarEventOption) error
}

// Mailer Mailer
type Mailer struct {
	MailProvider MailProvider
	Recipient    Recipient
	From       From
	FromDefault From
}

// NewMailer creating new mailer
func NewMailer(provider MailProvider, from From) *Mailer {
	return &Mailer{MailProvider: provider, From: from, FromDefault: from}
}


// Provider Set provider
func (m *Mailer) Provider(provider MailProvider) *Mailer {
	return NewMailer(provider, m.FromDefault)
}

func (m *Mailer) Reset() *Mailer {
	return NewMailer(m.MailProvider, m.FromDefault)
}

// Provider Set Sender
func (m *Mailer) Sender(from From) *Mailer {
	m.From = from
	return m
}

// AddTo Add to
func (m *Mailer) AddTo(to string) *Mailer {
	m.Recipient.To = append(m.Recipient.To, to)

	return m
}

// To Set to
func (m *Mailer) To(to ...string) *Mailer {
	m.Recipient.To = to

	return m
}

// AddCc Add cc
func (m *Mailer) AddCc(cc string) *Mailer {
	m.Recipient.Cc = append(m.Recipient.Cc, cc)

	return m
}

// Cc Set cc
func (m *Mailer) Cc(cc ...string) *Mailer {
	m.Recipient.Cc = cc

	return m
}

// AddBcc Add bcc
func (m *Mailer) AddBcc(bcc string) *Mailer {
	m.Recipient.Bcc = append(m.Recipient.Bcc, bcc)

	return m
}

// Bcc Set bcc
func (m *Mailer) Bcc(bcc ...string) *Mailer {
	m.Recipient.Bcc = bcc

	return m
}

// Send Send Email
func (m *Mailer) Send(mailable Mailable) error {
	mailContent := MailContent{
		Subject:  mailable.BuildSubject(),
		HtmlBody: mailable.BuildHTMLBody(),
		TextBody: mailable.BuildTextBody(),
	}
	return m.MailProvider.Send(m.Recipient, m.From, mailContent)
}

// Send Send Email
func (m *Mailer) SendWithCalendarEvent(mailable Mailable, options *CalendarEventOption) error {
	mailContent := MailContent{
		Subject:  mailable.BuildSubject(),
		HtmlBody: mailable.BuildHTMLBody(),
		TextBody: mailable.BuildTextBody(),
	}

	return m.MailProvider.SendWithCalendarEvent(m.Recipient, m.From, mailContent, options)
}
