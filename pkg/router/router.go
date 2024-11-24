package router

import (
	"github.com/RGaius/octopus/pkg/middleware"
	v1 "github.com/RGaius/octopus/pkg/router/v1"
	"github.com/RGaius/octopus/pkg/server/constant"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"os"
)

func InitRoutes(router *gin.Engine) {
	sessionSecret := "secret"
	if secretEnv := os.Getenv("SESSION_SECRET"); secretEnv != "" {
		sessionSecret = secretEnv
	}
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		HttpOnly: true,
		MaxAge:   constant.DefaultSessionExpiration,
		Path:     "/",
	})

	// use gin's crash free middleware
	router.Use(
		gin.Recovery(),
		requestid.New(),
		middleware.Logging(),
		gzip.Gzip(gzip.DefaultCompression),
		sessions.Sessions("cookies", store),
	)
	// v1 group
	v1Group := router.Group("/api/v1")
	v1.InitUserRouter(v1Group)
}
