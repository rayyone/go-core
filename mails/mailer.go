package mails

// Mailable represent mail interface
type Mailable interface {
	BuildSubject() string
	BuildHTMLBody() string
	BuildTextBody() string
}

type MailProvider interface {
	Send(recipient Recipient, subject string, htmlBody string, textBody string) error
	SendWithCalendarEvent(recipient Recipient, options *CalendarEventOption, subject string, htmlBody string, textBody string) error
}

// Mailer Mailer
type Mailer struct {
	MailProvider MailProvider
	Recipient    Recipient
}

// NewMailer creating new mailer
func NewMailer(provider MailProvider) *Mailer {
	return &Mailer{MailProvider: provider}
}

// Provider Set provider
func (m *Mailer) Provider(provider MailProvider) *Mailer {
	return NewMailer(provider)
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
	subject := mailable.BuildSubject()
	htmlBody := mailable.BuildHTMLBody()
	textBody := mailable.BuildTextBody()

	return m.MailProvider.Send(m.Recipient, subject, htmlBody, textBody)
}

// Send Send Email
func (m *Mailer) SendWithCalendarEvent(mailable Mailable, options *CalendarEventOption) error {
	subject := mailable.BuildSubject()
	htmlBody := mailable.BuildHTMLBody()
	textBody := mailable.BuildTextBody()

	return m.MailProvider.SendWithCalendarEvent(m.Recipient, options, subject, htmlBody, textBody)
}
