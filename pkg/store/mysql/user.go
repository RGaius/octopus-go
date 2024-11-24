package mysql

import (
	"github.com/RGaius/octopus/pkg/model/user"
)

type UserStore struct {
	Client *BaseDB
}

// GetUserByName returns a user object given a username
func (u *UserStore) GetUserByName(username string) (*user.User, error) {
	matchUser := new(user.User)
	exist, err := u.Client.Engine.Alias("user").Where("user.username=?", username).Get(matchUser)
	if !exist {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return matchUser, nil
}
