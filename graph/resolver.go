package graph

import "user-service/db/user"

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userRepository user.Repository
}

func NewResolver(userRepository user.Repository) *Resolver {
	return &Resolver{
		userRepository: userRepository,
	}
}
