package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"user-service/db/user"
	"user-service/kafka"
	"user-service/pkg"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var ErrCouldNotFindUser = errors.New("could not find user")

type Impl struct {
	repository user.Repository
}

func NewService(repository user.Repository) *Impl {
	return &Impl{
		repository: repository,
	}
}

func (s *Impl) GetUserById(ctx context.Context, log *zap.Logger, id uuid.UUID) (pkg.User, error) {
	dbUser, err := s.repository.GetUserById(ctx, id)
	if err != nil {
		log.Error("could not get user", zap.Error(err), zap.String("id", id.String()))
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.User{}, ErrCouldNotFindUser
		}

		return pkg.User{}, err
	}

	return MapUserToService(dbUser), nil
}

func (s *Impl) GetUsers(ctx context.Context, log *zap.Logger) ([]pkg.User, error) {
	dbUsers, err := s.repository.GetUsers(ctx)
	if err != nil {
		log.Error("could not get users", zap.Error(err))
		return nil, err
	}

	result := make([]pkg.User, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		result = append(result, MapUserToService(dbUser))
	}

	return result, nil
}

func (s *Impl) AddUser(ctx context.Context, log *zap.Logger, user pkg.User) (uuid.UUID, error) {
	id, err := s.repository.AddUser(ctx, MapUserToDb(user))
	if err != nil {
		log.Error("could not add user", zap.Error(err))
		return uuid.Nil, err
	}

	return id, nil
}

func (s *Impl) UpdateUser(ctx context.Context, log *zap.Logger, user pkg.User) error {
	err := s.repository.UpdateUser(ctx, MapUserToDb(user))
	if err != nil {
		log.Error("could not update user", zap.Error(err))
		return err
	}

	return nil
}

func (s *Impl) DeleteUser(ctx context.Context, log *zap.Logger, id uuid.UUID) error {
	err := s.repository.DeleteUser(ctx, id)
	if err != nil {
		log.Error("could not delete user", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	return nil
}

func (s *Impl) GetUserTicketsByUserId(ctx context.Context, log *zap.Logger, userId uuid.UUID) ([]pkg.UserTicket, error) {
	dbUserTickets, err := s.repository.GetUserTicketsByUserId(ctx, userId)
	if err != nil {
		log.Error("could not get user tickets", zap.Error(err))
		return nil, err
	}

	result := make([]pkg.UserTicket, 0, len(dbUserTickets))
	for _, dbUserTicket := range dbUserTickets {
		result = append(result, MapUserTicketToService(dbUserTicket))
	}

	return result, nil
}

func (s *Impl) CreateSubscriberForBookMessage(ctx context.Context, log *zap.Logger) kafka.Subscriber {
	return func(message kafka.Message, err error) {
		if err != nil {
			log.Error("could not create read message", zap.Error(err))
			return
		}

		var msg pkg.BookMessage
		err = json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Error("could not unmarshal message", zap.Error(err))
			return
		}

		err = s.repository.AddUserTicket(ctx, user.DbUserTicket{
			UserId:   msg.UserId,
			TicketId: msg.TicketId,
		})
		if err != nil {
			log.Error("could not add user ticket", zap.Error(err))
			return
		}
	}
}
