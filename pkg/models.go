package pkg

import "github.com/google/uuid"

type User struct {
	Id      uuid.UUID `json:"Id"`
	Email   string    `json:"Email"`
	Name    string    `json:"Name"`
	Surname string    `json:"Surname"`
}

type UserTicket struct {
	UserId   uuid.UUID `json:"UserId"`
	TicketId uuid.UUID `json:"TicketId"`
}

type BookMessage struct {
	UserId   uuid.UUID `json:"UserId"`
	TicketId uuid.UUID `json:"TicketId"`
}
