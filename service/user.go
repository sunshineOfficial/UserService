package service

import (
	"context"
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
}
