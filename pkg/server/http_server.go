package server

import (
	"fmt"
	"github.com/RGaius/octopus/pkg/router"
	"github.com/RGaius/octopus/pkg/server/constant"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"os"
)

type Authorizer interface {
	Authorize(req *http.Request) error
}

type HTTPServer struct {
	Router     *gin.Engine
	Server     *http.Server
	Authorizer *Authorizer
}

func (s *HTTPServer) RegisterRouter() error {
	router.InitRoutes(s.Router)
	return nil
}

func (s *HTTPServer) Run() error {
	s.Server.Handler = s.Router
	// address := fmt.Sprintf("%s:%d", constant.DefaultServerAddress, constant.DefaultServerPort)
	envPort := os.Getenv("LISTEN_PORT")
	var address string
	if envPort == "" {
		address = fmt.Sprintf("%s:%d", constant.DefaultServerHost, constant.DefaultServerPort)
	} else {
		address = fmt.Sprintf("%s:%s", constant.DefaultServerHost, envPort)
	}
	listener, err := net.Listen(constant.DefaultProtocol, address)
	if err != nil {
		return errors.Wrapf(err, "failed to listen address %s", address)
	}
	err = s.Server.Serve(listener)
	if err != nil {
		return errors.Wrap(err, "failed to start server")
	}
	return nil
}

func NewHTTPServer() *HTTPServer {
	if os.Getenv("DEBUG") != "true" {
		gin.SetMode(gin.ReleaseMode)
	}
	return &HTTPServer{
		Router: gin.New(),
		Server: &http.Server{},
	}
}
