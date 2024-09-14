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

func GetUserByIdHandler(userService service.User, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := chi.URLParam(r, "id")
		id, err := uuid.Parse(idRaw)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, "wrong id")
			return
		}

		result, err := userService.GetUserById(r.Context(), log, id)
		if err != nil {
			if errors.Is(err, user.ErrCouldNotFindUser) {
				render.Status(r, http.StatusNoContent)
				render.PlainText(w, r, err.Error())
				return
			}

			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, result)
		return
	}
}

func GetUsersHandler(userService service.User, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := userService.GetUsers(r.Context(), log)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, result)
		return
	}
}

func AddUserHandler(userService service.User, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u pkg.User
		err := render.DecodeJSON(r.Body, &u)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		id, err := userService.AddUser(r.Context(), log, u)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.PlainText(w, r, id.String())
		return
	}
}

func UpdateUserHandler(userService service.User, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := chi.URLParam(r, "id")
		id, err := uuid.Parse(idRaw)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, "wrong id")
			return
		}

		var u pkg.User
		err = render.DecodeJSON(r.Body, &u)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		u.Id = id

		err = userService.UpdateUser(r.Context(), log, u)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.PlainText(w, r, "ok")
		return
	}
}

func DeleteUserHandler(userService service.User, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := chi.URLParam(r, "id")
		id, err := uuid.Parse(idRaw)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, "wrong id")
			return
		}

		err = userService.DeleteUser(r.Context(), log, id)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		render.Status(r, http.StatusOK)
		render.PlainText(w, r, "ok")
		return
	}
}
