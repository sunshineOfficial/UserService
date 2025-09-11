package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"user-service/pkg"
	"user-service/service"
	"user-service/service/user"

	"github.com/google/uuid"
	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

type idVars struct {
	id string
}

// GetUserByIdHandler получает пользователя по ID
//
//	@Summary	Получает пользователя по ID
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	pkg.User
//	@Success	204
//	@Failure	400
//	@Router		/user/{id} [get]
func GetUserByIdHandler(userService service.User) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars idVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("parse vars: %w", err)
		}

		id, err := uuid.Parse(vars.id)
		if err != nil {
			return fmt.Errorf("parse id: %w", err)
		}

		result, err := userService.GetUserById(c.Ctx(), c.Log(), id)
		if err != nil {
			if errors.Is(err, user.ErrCouldNotFindUser) {
				c.Write(http.StatusNoContent)
				return nil
			}

			return fmt.Errorf("get user by id: %w", err)
		}

		return c.WriteJson(http.StatusOK, result)
	}
}

// GetUsersHandler получает всех пользователей по ID
//
//	@Summary	Получает всех пользователей по ID
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]pkg.User
//	@Failure	400
//	@Router		/user [get]
func GetUsersHandler(userService service.User) gorouter.Handler {
	return func(c gorouter.Context) error {
		result, err := userService.GetUsers(c.Ctx(), c.Log())
		if err != nil {
			return fmt.Errorf("get users: %w", err)
		}

		return c.WriteJson(http.StatusOK, result)
	}
}

// AddUserHandler добавляет нового пользователя
//
//	@Summary	Добавляет нового пользователя
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		user	body		pkg.User	true	"User"
//	@Success	200		{object}	string
//	@Failure	400
//	@Router		/user [post]
func AddUserHandler(userService service.User) gorouter.Handler {
	return func(c gorouter.Context) error {
		var u pkg.User
		if err := c.ReadJson(&u); err != nil {
			return fmt.Errorf("parse json: %w", err)
		}

		id, err := userService.AddUser(c.Ctx(), c.Log(), u)
		if err != nil {
			return fmt.Errorf("add user: %w", err)
		}

		return c.WriteJson(http.StatusOK, id.String())
	}
}

// UpdateUserHandler обновляет пользователя
//
//	@Summary	Обновляет пользователя
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string		true	"User ID"
//	@Param		user	body	pkg.User	true	"User"
//	@Success	200
//	@Failure	400
//	@Router		/user [put]
func UpdateUserHandler(userService service.User) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars idVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("parse vars: %w", err)
		}

		id, err := uuid.Parse(vars.id)
		if err != nil {
			return fmt.Errorf("parse id: %w", err)
		}

		var u pkg.User
		if err = c.ReadJson(&u); err != nil {
			return fmt.Errorf("parse json: %w", err)
		}

		u.Id = id

		err = userService.UpdateUser(c.Ctx(), c.Log(), u)
		if err != nil {
			return fmt.Errorf("update user: %w", err)
		}

		c.Write(http.StatusOK)
		return nil
	}
}

// DeleteUserHandler удаляет пользователя по ID
//
//	@Summary	Удаляет пользователя по ID
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"User ID"
//	@Success	200
//	@Failure	400
//	@Router		/user/{id} [delete]
func DeleteUserHandler(userService service.User) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars idVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("parse vars: %w", err)
		}

		id, err := uuid.Parse(vars.id)
		if err != nil {
			return fmt.Errorf("parse id: %w", err)
		}

		err = userService.DeleteUser(c.Ctx(), c.Log(), id)
		if err != nil {
			return fmt.Errorf("delete user: %w", err)
		}

		c.Write(http.StatusOK)
		return nil
	}
}

// GetUserTicketsByUserIdHandler получает билеты пользователя по его ID
//
//	@Summary	Получает билеты пользователя по его ID
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	[]pkg.UserTicket
//	@Failure	400
//	@Router		/user/{id}/tickets [get]
func GetUserTicketsByUserIdHandler(userService service.User) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars idVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("parse vars: %w", err)
		}

		id, err := uuid.Parse(vars.id)
		if err != nil {
			return fmt.Errorf("parse id: %w", err)
		}

		result, err := userService.GetUserTicketsByUserId(c.Ctx(), c.Log(), id)
		if err != nil {
			return fmt.Errorf("get user tickets: %w", err)
		}

		return c.WriteJson(http.StatusOK, result)
	}
}
