package mails

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"strings"
	"time"

	loghelper "github.com/rayyone/go-core/helpers/log"
	"github.com/rayyone/go-core/ryerr"
)

// SMTPMailer SMTP SMTPMailer
type SMTPMailer struct {
	config SMTPConfiguration
}

type SMTPConfiguration struct {
	Host     string
	Port     int
	Email    string
	Password string
	Name     string
}

func (m *SMTPMailer) Send(recipient Recipient, subject string, htmlBody string, textBody string) error {
	smtpAuth := smtp.PlainAuth(m.config.Name, m.config.Email, m.config.Password, m.config.Host)
	smtpAddr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)

	msg := m.buildMessage(recipient, subject, htmlBody, textBody)

	start := time.Now()
	loghelper.PrintYellowf(
		"[SMTP] Sending email to: %s, cc: %s, bcc: %s",
		strings.Join(recipient.To, ", "),
		strings.Join(recipient.Cc, ", "),
		strings.Join(recipient.Bcc, ", "),
	)
	err := smtp.SendMail(smtpAddr, smtpAuth, m.config.Email, recipient.To, []byte(msg))
	if err != nil {
		loghelper.PrintRedf("[SMTP] Send email completed with error in %.2fs", time.Since(start).Seconds())
		return ryerr.NewAndDontReport(fmt.Sprintf("Error: Cannot send email via SMTP provider. Error: %v", err))
	}
	loghelper.PrintYellowf("[SMTP] Send email completed in %.2fs", time.Since(start).Seconds())

	return nil
}

func (m *SMTPMailer) SendWithCalendarEvent(recipient Recipient, options *CalendarEventOption, subject string, htmlBody string, textBody string) error {
	smtpAuth := smtp.PlainAuth(m.config.Name, m.config.Email, m.config.Password, m.config.Host)
	smtpAddr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)

	msg := m.buildCalendarInvitationMessage(recipient, options, subject, htmlBody, textBody)

	err := smtp.SendMail(smtpAddr, smtpAuth, m.config.Email, recipient.To, []byte(msg))
	if err != nil {
		return ryerr.NewAndDontReport(fmt.Sprintf("Error: Cannot send email via SMTP provider. Error: %v", err))
	}

	return nil
}

// NewSMTPMailer New smtp mailer
func NewSMTPMailer(conf SMTPConfiguration) *SMTPMailer {
	return &SMTPMailer{config: conf}
}

func (m *SMTPMailer) buildMessage(recipient Recipient, subject string, htmlBody string, textBody string) string {
	writer := multipart.NewWriter(bytes.NewBufferString(""))

	msg := fmt.Sprintf("From: %s <%s>\r\n", m.config.Name, m.config.Email)
	if len(recipient.To) > 0 {
		msg += fmt.Sprintf("To: %s\r\n", strings.Join(recipient.To, ";"))
	}
	if len(recipient.Cc) > 0 {
		msg += fmt.Sprintf("Cc: %s\r\n", strings.Join(recipient.Cc, ";"))
	}
	if len(recipient.Bcc) > 0 {
		msg += fmt.Sprintf("Bcc: %s\r\n", strings.Join(recipient.Bcc, ";"))
	}
	msg += "Subject: " + subject + "\r\n"
	msg += "MIME-version: 1.0;"
	msg += getAlternativeMultipartStart(writer)
	msg += getContentTypeWithBoundary(writer, "text/plain", "UTF-8", "8bit")
	msg += "\r\n" + textBody
	msg += getContentTypeWithBoundary(writer, "text/html", "UTF-8", "8bit")
	msg += "\r\n" + htmlBody
	msg += getMultipartBoundaryEnd(writer)

	return msg
}

func (m *SMTPMailer) buildCalendarInvitationMessage(recipient Recipient, options *CalendarEventOption, subject string, htmlBody string, textBody string) string {
	mixedBoundaryWriter := multipart.NewWriter(bytes.NewBufferString(""))
	alternativeBoundaryWriter := multipart.NewWriter(bytes.NewBufferString(""))

	msg := fmt.Sprintf("From: %s <%s>\r\n", m.config.Name, m.config.Email)
	if len(recipient.To) > 0 {
		msg += fmt.Sprintf("To: %s\r\n", strings.Join(recipient.To, ";"))
	}
	if len(recipient.Cc) > 0 {
		msg += fmt.Sprintf("Cc: %s\r\n", strings.Join(recipient.Cc, ";"))
	}
	if len(recipient.Bcc) > 0 {
		msg += fmt.Sprintf("Bcc: %s\r\n", strings.Join(recipient.Bcc, ";"))
	}
	msg += "Subject: " + subject + "\r\n"
	msg += "MIME-version: 1.0;"
	msg += getMixedMultipartStart(mixedBoundaryWriter)
	msg += getMultipartBoundaryOpen(mixedBoundaryWriter)
	msg += getAlternativeMultipartStart(alternativeBoundaryWriter)
	msg += getContentTypeWithBoundary(alternativeBoundaryWriter, "text/plain", "UTF-8", "8bit")
	msg += "\r\n" + textBody
	msg += getContentTypeWithBoundary(alternativeBoundaryWriter, "text/html", "UTF-8", "8bit")
	msg += "\r\n" + htmlBody
	msg += getContentTypeWithBoundary(alternativeBoundaryWriter, "text/calendar; method=REQUEST", "UTF-8", "7bit")
	msg += getCalendarBody(options)
	msg += getMultipartBoundaryEnd(alternativeBoundaryWriter)
	msg += getICSAttachmentWithBoundary(mixedBoundaryWriter)
	msg += getMultipartBoundaryEnd(mixedBoundaryWriter)

	return msg
}

func getMixedMultipartStart(writer *multipart.Writer) string {
	content := `
Content-Type: multipart/mixed; charset="UTF-8"; boundary="%s"
`
	return fmt.Sprintf(content, writer.Boundary())
}

func getAlternativeMultipartStart(writer *multipart.Writer) string {
	content := `
Content-Type: multipart/alternative; charset="UTF-8"; boundary="%s"
`
	return fmt.Sprintf(content, writer.Boundary())
}

func getMultipartBoundaryOpen(writer *multipart.Writer) string {
	boundary := "\n--%s"
	return fmt.Sprintf(boundary, writer.Boundary())
}

func getContentTypeWithBoundary(writer *multipart.Writer, contentType string, charset string, encoding string) string {
	contentTypeFormat := `
--%s
Content-Type: %s; charset="%s";
Content-Transfer-Encoding: %s
`
	return fmt.Sprintf(contentTypeFormat, writer.Boundary(), contentType, charset, encoding)
}

func getCalendarBody(options *CalendarEventOption) string {
	var attendee string
	for _, a := range options.Attendees {
		attendee += fmt.Sprintf("\nATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=%s;RSVP=TRUE;CN=%s;X-NUM-GUESTS=0:mailto:%s", a.Status, a.Name, a.Email)
	}

	if options.AppointmentMethod == "" {
		options.AppointmentMethod = APPOINTMENT_TYPE_REQUEST
	}

	body := `
BEGIN:VCALENDAR
PRODID:%s
METHOD:%s
VERSION:2.0
BEGIN:VEVENT
UID:%s
SEQUENCE:%d
STATUS:%s
ORGANIZER:mailto:%s
%s
DTSTAMP:%s
DTSTART:%s
DTEND:%s
SUMMARY:%s
DESCRIPTION:%s
END:VEVENT
END:VCALENDAR
`
	iso8601 := "20060102T150405Z"
	return fmt.Sprintf(body,
		options.ProdID,
		options.AppointmentMethod,
		options.EventID,
		options.Sequence,
		options.Status,
		options.Organizer,
		attendee,
		time.Now().UTC().Format(iso8601),
		options.StartDateTime.UTC().Format(iso8601),
		options.EndDateTime.UTC().Format(iso8601),
		options.Summary,
		options.Description,
	)
}

func getICSAttachmentWithBoundary(writer *multipart.Writer) string {
	content := `
--%s
Content-Type: application/ics; name="invite.ics"
Content-Disposition: attachment; filename="invite.ics"
Content-Transfer-Encoding: base64
`
	return fmt.Sprintf(content, writer.Boundary())
}

func getMultipartBoundaryEnd(writer *multipart.Writer) string {
	content := `
--%s--
`
	return fmt.Sprintf(content, writer.Boundary())
}
