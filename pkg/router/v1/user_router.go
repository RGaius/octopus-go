package v1

import (
	h "github.com/RGaius/octopus/pkg/handler"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(g *gin.RouterGroup) {
	g.POST("/login", h.Wrap(h.Login))
	g.POST("/logout", h.Wrap(h.Logout))
}
