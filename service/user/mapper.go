package user

import (
	"user-service/db/user"
	"user-service/pkg"
)

func MapUserToService(db user.DbUser) pkg.User {
	return pkg.User{
		Id:      db.Id,
		Email:   db.Email,
		Name:    db.Name,
		Surname: db.Surname,
	}
}

func MapUserToDb(service pkg.User) user.DbUser {
	return user.DbUser{
		Id:      service.Id,
		Email:   service.Email,
		Name:    service.Name,
		Surname: service.Surname,
	}
}

func MapUserTicketToService(db user.DbUserTicket) pkg.UserTicket {
	return pkg.UserTicket{
		UserId:   db.UserId,
		TicketId: db.TicketId,
	}
}
