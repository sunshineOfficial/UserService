package service

import (
	"context"
	"user-service/kafka"
	"user-service/pkg"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type User interface {
	GetUserById(ctx context.Context, log *zap.Logger, id uuid.UUID) (pkg.User, error)
	GetUsers(ctx context.Context, log *zap.Logger) ([]pkg.User, error)
	AddUser(ctx context.Context, log *zap.Logger, user pkg.User) (uuid.UUID, error)
	UpdateUser(ctx context.Context, log *zap.Logger, user pkg.User) error
	DeleteUser(ctx context.Context, log *zap.Logger, id uuid.UUID) error
	GetUserTicketsByUserId(ctx context.Context, log *zap.Logger, userId uuid.UUID) ([]pkg.UserTicket, error)
	CreateSubscriberForBookMessage(ctx context.Context, log *zap.Logger) kafka.Subscriber
}
