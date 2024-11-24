package handler

import (
	"github.com/RGaius/octopus/pkg/errors"
	"github.com/RGaius/octopus/pkg/model/param"
	"github.com/RGaius/octopus/pkg/model/user"
	"github.com/RGaius/octopus/pkg/store"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) (*user.Account, error) {
	// 获取登录参数
	loginParam := &param.LoginParam{}
	err := c.BindJSON(loginParam)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}
	defaultStore, err := store.Load()
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}

	// 通过用户名获取用户信息
	matchUser, err := defaultStore.GetUserByName(loginParam.Username)
	if err != nil {
		return nil, errors.NewNotFound(err.Error())
	}
	if matchUser == nil {
		return nil, errors.NewNotFound("当前用户不存在")
	}
	return nil, nil
}

func Logout(c *gin.Context) (string, error) {
	return "", nil
}
