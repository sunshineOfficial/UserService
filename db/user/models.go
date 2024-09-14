package user

import "github.com/google/uuid"

type DbUser struct {
	Id      uuid.UUID `db:"id"`
	Email   string    `db:"email"`
	Name    string    `db:"name"`
	Surname string    `db:"surname"`
}
