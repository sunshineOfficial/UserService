package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"user-service/db/user"
	"user-service/pkg"

	"github.com/google/uuid"
	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
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

func (s *Impl) GetUserById(ctx goctx.Context, log golog.Logger, id uuid.UUID) (pkg.User, error) {
	dbUser, err := s.repository.GetUserById(ctx, id)
	if err != nil {
		log.Errorf("could not get user %s: %v", id, err)
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.User{}, ErrCouldNotFindUser
		}

		return pkg.User{}, err
	}

	return MapUserToService(dbUser), nil
}

func (s *Impl) GetUsers(ctx goctx.Context, log golog.Logger) ([]pkg.User, error) {
	dbUsers, err := s.repository.GetUsers(ctx)
	if err != nil {
		log.Errorf("could not get users: %v", err)
		return nil, err
	}

	result := make([]pkg.User, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		result = append(result, MapUserToService(dbUser))
	}

	return result, nil
}

func (s *Impl) AddUser(ctx goctx.Context, log golog.Logger, user pkg.User) (uuid.UUID, error) {
	id, err := s.repository.AddUser(ctx, MapUserToDb(user))
	if err != nil {
		log.Errorf("could not add user: %v", err)
		return uuid.Nil, err
	}

	return id, nil
}

func (s *Impl) UpdateUser(ctx goctx.Context, log golog.Logger, user pkg.User) error {
	err := s.repository.UpdateUser(ctx, MapUserToDb(user))
	if err != nil {
		log.Errorf("could not update user: %v", err)
		return err
	}

	return nil
}

func (s *Impl) DeleteUser(ctx goctx.Context, log golog.Logger, id uuid.UUID) error {
	err := s.repository.DeleteUser(ctx, id)
	if err != nil {
		log.Errorf("could not delete user %s: %v", id, err)
		return err
	}

	return nil
}

func (s *Impl) GetUserTicketsByUserId(ctx goctx.Context, log golog.Logger, userId uuid.UUID) ([]pkg.UserTicket, error) {
	dbUserTickets, err := s.repository.GetUserTicketsByUserId(ctx, userId)
	if err != nil {
		log.Errorf("could not get user tickets: %v", err)
		return nil, err
	}

	result := make([]pkg.UserTicket, 0, len(dbUserTickets))
	for _, dbUserTicket := range dbUserTickets {
		result = append(result, MapUserTicketToService(dbUserTicket))
	}

	return result, nil
}

func (s *Impl) CreateSubscriberForBookMessage(ctx context.Context, log golog.Logger) gokafka.Subscriber {
	return func(message gokafka.Message, err error) {
		if err != nil {
			log.Errorf("could not create read message: %v", err)
			return
		}

		var msg pkg.BookMessage
		err = json.Unmarshal(message.Value, &msg)
		if err != nil {
			log.Errorf("could not unmarshal message: %v", err)
			return
		}

		err = s.repository.AddUserTicket(ctx, user.DbUserTicket{
			UserId:   msg.UserId,
			TicketId: msg.TicketId,
		})
		if err != nil {
			log.Errorf("could not add user ticket: %v", err)
			return
		}

		log.Debug(fmt.Sprintf("consumed book message: %v", msg))
	}
}
