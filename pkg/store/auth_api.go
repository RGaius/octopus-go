package store

import "github.com/RGaius/octopus/pkg/model/user"

type AuthStore interface {
	UserStore
}
type UserStore interface {
	// GetUserByName returns the user with the given username.
	GetUserByName(username string) (*user.User, error)
}
