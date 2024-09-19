package user

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Impl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Impl {
	return Impl{
		db: db,
	}
}

//go:embed sql/get_user_by_id.sql
var getUserByIdSql string

func (r Impl) GetUserById(ctx context.Context, id uuid.UUID) (DbUser, error) {
	var user DbUser
	err := r.db.GetContext(ctx, &user, getUserByIdSql, id)

	return user, err
}

//go:embed sql/get_users.sql
var getUsersSql string

func (r Impl) GetUsers(ctx context.Context) ([]DbUser, error) {
	users := make([]DbUser, 0)

	err := r.db.SelectContext(ctx, &users, getUsersSql)
	if errors.Is(err, sql.ErrNoRows) {
		return users, nil
	}

	return users, err
}

//go:embed sql/add_user.sql
var addUserSql string

func (r Impl) AddUser(ctx context.Context, user DbUser) (id uuid.UUID, err error) {
	rows, err := r.db.NamedQueryContext(ctx, addUserSql, user)
	if err != nil {
		return uuid.Nil, err
	}

	defer func(rows *sqlx.Rows) {
		if tempErr := rows.Close(); tempErr != nil {
			err = tempErr
		}
	}(rows)

	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return uuid.Nil, err
		}
	}

	return
}

//go:embed sql/update_user.sql
var updateUserSql string

func (r Impl) UpdateUser(ctx context.Context, user DbUser) error {
	_, err := r.db.NamedExecContext(ctx, updateUserSql, user)

	return err
}

//go:embed sql/delete_user.sql
var deleteUserSql string

func (r Impl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, deleteUserSql, id)

	return err
}

//go:embed sql/get_user_tickets_by_user_id.sql
var getUserTicketsByUserIdSql string

func (r Impl) GetUserTicketsByUserId(ctx context.Context, userId uuid.UUID) ([]DbUserTicket, error) {
	userTickets := make([]DbUserTicket, 0)

	err := r.db.SelectContext(ctx, &userTickets, getUserTicketsByUserIdSql, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return userTickets, nil
	}

	return userTickets, err
}

//go:embed sql/add_user_ticket.sql
var addUserTicketSql string

func (r Impl) AddUserTicket(ctx context.Context, userTicket DbUserTicket) error {
	_, err := r.db.NamedExecContext(ctx, addUserTicketSql, userTicket)

	return err
}
