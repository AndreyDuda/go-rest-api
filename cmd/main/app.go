package main

import (
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"rest_api/internal/user"
	"rest_api/pkg/logging"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()
	handler := user.NewHandler(logger)

	logger.Info("register user handler")
	handler.Register(router)

	start(router)
}

func start(router *httprouter.Router) {
	logger := logging.GetLogger()
	logger.Info("start application")
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("server are listening port 0.0.0.0:1234")
	logger.Fatal(server.Serve(listener))
}
