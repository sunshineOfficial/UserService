package service

import (
	"context"
	"user-service/pkg"

	"github.com/google/uuid"
	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
)

type User interface {
	GetUserById(ctx goctx.Context, log golog.Logger, id uuid.UUID) (pkg.User, error)
	GetUsers(ctx goctx.Context, log golog.Logger) ([]pkg.User, error)
	AddUser(ctx goctx.Context, log golog.Logger, user pkg.User) (uuid.UUID, error)
	UpdateUser(ctx goctx.Context, log golog.Logger, user pkg.User) error
	DeleteUser(ctx goctx.Context, log golog.Logger, id uuid.UUID) error
	GetUserTicketsByUserId(ctx goctx.Context, log golog.Logger, userId uuid.UUID) ([]pkg.UserTicket, error)
	CreateSubscriberForBookMessage(ctx context.Context, log golog.Logger) gokafka.Subscriber
}
