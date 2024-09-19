package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (DbUser, error)
	GetUsers(ctx context.Context) ([]DbUser, error)
	AddUser(ctx context.Context, user DbUser) (uuid.UUID, error)
	UpdateUser(ctx context.Context, user DbUser) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUserTicketsByUserId(ctx context.Context, userId uuid.UUID) ([]DbUserTicket, error)
	AddUserTicket(ctx context.Context, userTicket DbUserTicket) error
}
