package mails

type Mailable interface {
	Subject() string
	HTMLBody() string
	TextBody() string
	Header() map[string]string
	Attachments() []Attachment
}
type From struct {
	Address string
	Name    string
}
type Envelope struct {
	To   []string
	Cc   []string
	Bcc  []string
	From *From
}
type Message struct {
	Envelope
	Subject     string
	Text        string
	HTML        string
	Attachments []Attachment
	Header      map[string]string
}

type MailProvider interface {
	Send(content Message) error
}

type Mailer struct {
	MailProvider MailProvider
}

type MailBuilder struct {
	Mailer
	Message Message
}

var mailer *Mailer

func NewMailer(provider MailProvider) *Mailer {
	mailer = &Mailer{MailProvider: provider}
	return mailer
}
func (m *Mailer) Send(msg Message) error {
	return m.MailProvider.Send(msg)
}

func (m *Mailer) SendMailable(envelope Envelope, mailable Mailable) error {
	msg := Message{Envelope: envelope}
	msg.HTML = mailable.HTMLBody()
	msg.Text = mailable.TextBody()
	msg.Subject = mailable.Subject()
	msg.Header = mailable.Header()
	msg.Attachments = mailable.Attachments()
	return m.MailProvider.Send(msg)
}
func (m *Mailer) MailBuilder() *MailBuilder {
	var mailBuilder MailBuilder
	mailBuilder.MailProvider = m.MailProvider
	return &mailBuilder
}
func (m *MailBuilder) To(to string) *MailBuilder {
	if m.Message.To == nil {
		m.Message.To = []string{to}
	} else {
		m.Message.To = append(m.Message.To, to)
	}
	return m
}
func (m *MailBuilder) From(from From) *MailBuilder {
	m.Message.From = &from
	return m
}
func (m *MailBuilder) Bcc(bcc string) *MailBuilder {
	if m.Message.Bcc == nil {
		m.Message.Bcc = []string{bcc}
	} else {
		m.Message.Bcc = append(m.Message.Bcc, bcc)
	}
	return m
}
func (m *MailBuilder) Cc(cc string) *MailBuilder {
	if m.Message.Cc == nil {
		m.Message.Cc = []string{cc}
	} else {
		m.Message.Cc = append(m.Message.Cc, cc)
	}
	return m
}
func (m *MailBuilder) Subject(subject string) *MailBuilder {
	m.Message.Subject = subject
	return m
}
func (m *MailBuilder) HTML(html string) *MailBuilder {
	m.Message.HTML = html
	return m
}
func (m *MailBuilder) Text(text string) *MailBuilder {
	m.Message.Text = text
	return m
}
func (m *MailBuilder) Header(header map[string]string) *MailBuilder {
	m.Message.Header = header
	return m
}
func (m *MailBuilder) Attachment(attachment Attachment) *MailBuilder {
	if m.Message.Attachments == nil {
		m.Message.Attachments = []Attachment{attachment}
	} else {
		m.Message.Attachments = append(m.Message.Attachments, attachment)
	}
	return m
}
func (m *MailBuilder) Deliver() error {
	return m.Send(m.Message)
}
