package user

import (
	"user-service/db/user"
	"user-service/pkg"
)

func MapToService(db user.DbUser) pkg.User {
	return pkg.User{
		Id:      db.Id,
		Email:   db.Email,
		Name:    db.Name,
		Surname: db.Surname,
	}
}

func MapToDb(service pkg.User) user.DbUser {
	return user.DbUser{
		Id:      service.Id,
		Email:   service.Email,
		Name:    service.Name,
		Surname: service.Surname,
	}
}
