package pkg

import "github.com/google/uuid"

type User struct {
	Id      uuid.UUID `json:"Id"`
	Email   string    `json:"Email"`
	Name    string    `json:"Name"`
	Surname string    `json:"Surname"`
}
