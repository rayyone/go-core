package mails

// Mailable represent mail interface
type Mailable interface {
	BuildSubject() string
	BuildHTMLBody() string
	BuildTextBody() string
	BuildHeader() map[string]string
	CallbackFunc() func(args ...interface{})
}

type MailProvider interface {
	Send(content MailContent) error
	SendWithCalendarEvent(content MailContent, options *CalendarEventOption) error
}

// Mailer Mailer
type Mailer struct {
	MailProvider MailProvider
	From         From
}
var mailer *Mailer
func GetCurrentMailer() *Mailer {
	return mailer
}
// NewMailer creating new mailer
func NewMailer(provider MailProvider, from From) *Mailer {
	mailer = &Mailer{MailProvider: provider, From: from}
	return mailer
}

// Provider Set provider
func (m *Mailer) Provider(provider MailProvider) *Mailer {
	return NewMailer(provider, m.From)
}

func (m *Mailer) GetContent(mailable Mailable) MailContent {
	return MailContent{
		Subject:      mailable.BuildSubject(),
		HtmlBody:     mailable.BuildHTMLBody(),
		TextBody:     mailable.BuildTextBody(),
		Header:       mailable.BuildHeader(),
		CallbackFunc: mailable.CallbackFunc(),
		From: m.From,
	}
}

// Send Send Email
func (m *Mailer) Send(mailContent MailContent) error {
	return m.MailProvider.Send(mailContent)
}

// Send Send Email
func (m *Mailer) SendWithCalendarEvent(mailContent MailContent, options *CalendarEventOption) error {

	return m.MailProvider.SendWithCalendarEvent(mailContent, options)
}

func (c *MailContent) SetFrom(from From) *MailContent {
	c.From = from
	return c
}

// AddTo Add to
func (c MailContent) AddTo(to string) MailContent {
	c.Recipient.To = append(c.Recipient.To, to)

	return c
}

// To Set to
func (c MailContent) To(to ...string) MailContent {
	c.Recipient.To = to

	return c
}

// AddCc Add cc
func (c MailContent) AddCc(cc string) MailContent {
	c.Recipient.Cc = append(c.Recipient.Cc, cc)

	return c
}

// Cc Set cc
func (c MailContent) Cc(cc ...string) MailContent {
	c.Recipient.Cc = cc

	return c
}

// AddBcc Add bcc
func (c MailContent) AddBcc(bcc string) MailContent {
	c.Recipient.Bcc = append(c.Recipient.Bcc, bcc)

	return c
}

// Bcc Set bcc
func (c MailContent) Bcc(bcc ...string) MailContent {
	c.Recipient.Bcc = bcc
	return c
}
