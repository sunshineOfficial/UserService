package handlers

import (
	"errors"
	"net/http"
	"user-service/pkg"
	"user-service/service"
	"user-service/service/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GetUserByIdHandler получает пользователя по ID
//
//	@Summary	Получает пользователя по ID
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	pkg.User
//	@Success	204	{object}	string
//	@Failure	400	{object}	string
//	@Router		/user/{id} [get]
func GetUserByIdHandler(userService service.User, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := chi.URLParam(r, "id")
		id, err := uuid.Parse(idRaw)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "wrong id")
			return
		}

		result, err := userService.GetUserById(r.Context(), log, id)
		if err != nil {
			if errors.Is(err, user.ErrCouldNotFindUser) {
				render.Status(r, http.StatusNoContent)
				render.JSON(w, r, err.Error())
				return
			}

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, result)
		return
	}
}

// GetUsersHandler получает всех пользователей по ID
//
//	@Summary	Получает всех пользователей по ID
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]pkg.User
//	@Failure	400	{object}	string
//	@Router		/user [get]
func GetUsersHandler(userService service.User, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := userService.GetUsers(r.Context(), log)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, result)
		return
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
//	@Failure	400		{object}	string
//	@Router		/user [post]
func AddUserHandler(userService service.User, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u pkg.User
		err := render.DecodeJSON(r.Body, &u)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err.Error())
			return
		}

		id, err := userService.AddUser(r.Context(), log, u)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, id.String())
		return
	}
}

// UpdateUserHandler обновляет пользователя
//
//	@Summary	Обновляет пользователя
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string		true	"User ID"
//	@Param		user	body		pkg.User	true	"User"
//	@Success	200		{object}	string
//	@Failure	400		{object}	string
//	@Router		/user [put]
func UpdateUserHandler(userService service.User, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := chi.URLParam(r, "id")
		id, err := uuid.Parse(idRaw)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "wrong id")
			return
		}

		var u pkg.User
		err = render.DecodeJSON(r.Body, &u)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err.Error())
			return
		}

		u.Id = id

		err = userService.UpdateUser(r.Context(), log, u)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, "ok")
		return
	}
}

// DeleteUserHandler удаляет пользователя по ID
//
//	@Summary	Удаляет пользователя по ID
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	string
//	@Failure	400	{object}	string
//	@Router		/user/{id} [delete]
func DeleteUserHandler(userService service.User, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := chi.URLParam(r, "id")
		id, err := uuid.Parse(idRaw)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "wrong id")
			return
		}

		err = userService.DeleteUser(r.Context(), log, id)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, "ok")
		return
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
//	@Failure	400	{object}	string
//	@Router		/user/{id}/tickets [get]
func GetUserTicketsByUserIdHandler(userService service.User, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := chi.URLParam(r, "id")
		id, err := uuid.Parse(idRaw)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "wrong id")
			return
		}

		result, err := userService.GetUserTicketsByUserId(r.Context(), log, id)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, result)
		return
	}
}
