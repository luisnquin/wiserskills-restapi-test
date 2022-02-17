package models

import "time"

type (
	Event struct {
		Id         uint16    `json:"id" sql:"id,pk"`
		Name       string    `json:"name" sql:"name"`
		Created_at time.Time `json:"created_at" sql:"created_at"`
	}
	Events []Event
)

type (
	Participant struct {
		Id        uint64 `json:"id" sql:"id,pk"`
		Firstname string `json:"firstname" sql:"firstname"`
		Lastname  string `json:"lastname" sql:"lastname"`
		Age       uint8  `json:"age" sql:"age"`
	}
	Participants []Participant
)

type (
	Ticket struct {
		Id          uint32 `json:"id" sql:"id,pk"`
		Participant uint64 `json:"participant" sql:"participant"`
		Event       uint16 `json:"event" sql:"event"`
	}
	Tickets []Ticket

	TicketView struct {
		Id          uint32 `json:"id" sql:"id,pk"`
		Participant string `json:"participant" sql:"participant"`
		Event       string `json:"event" sql:"event"`
	}
	TicketViews []TicketView
)
