package mails

import "time"

// Recipient Recipient
type Recipient struct {
	To  []string
	Cc  []string
	Bcc []string
}

type AttendeeStatus string

const (
	STATUS_ACCEPTED     AttendeeStatus = "ACCEPTED"
	STATUS_NEEDS_ACTION AttendeeStatus = "NEEDS-ACTION"
)

type EventStatus string

const (
	EVENT_STATUS_TENTATIVE EventStatus = "TENTATIVE"
	EVENT_STATUS_CONFIRMED EventStatus = "CONFIRMED"
	EVENT_STATUS_CANCELLED EventStatus = "CANCELLED"
)

type AppointmentMethod string

const (
	APPOINTMENT_TYPE_REQUEST AppointmentMethod = "REQUEST"
	APPOINTMENT_TYPE_PUBLISH AppointmentMethod = "PUBLISH"
	APPOINTMENT_TYPE_CANCEL  AppointmentMethod = "CANCEL"
)

type Attendee struct {
	Name   string
	Email  string
	Status AttendeeStatus
}

type CalendarEventOption struct {
	ProdID            string
	AppointmentMethod AppointmentMethod
	EventID           string
	Sequence          int
	Status            EventStatus
	Organizer         string
	Summary           string
	Description       string
	StartDateTime     time.Time
	EndDateTime       time.Time
	Attendees         []Attendee
}

func (opt *CalendarEventOption) AddAttendee(attendee Attendee) *CalendarEventOption {
	opt.Attendees = append(opt.Attendees, attendee)
	return opt
}
