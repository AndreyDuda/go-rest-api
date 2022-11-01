package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest_api/internal/config"
	"rest_api/internal/user"
	"rest_api/internal/user/db"
	"rest_api/pkg/client/mongodb"
	"rest_api/pkg/logging"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()
	handler := user.NewHandler(logger)

	cfg := config.GetConfig()

	cfgMongo := cfg.Mongodb
	mongoDBClient, err := mongodb.NewClient(
		context.Background(),
		cfgMongo.Host,
		cfgMongo.Port,
		cfgMongo.Username,
		cfgMongo.Password,
		cfgMongo.Database,
		cfgMongo.AuthDB,
	)
	if err != nil {
		panic(err)
	}

	user1 := user.User{
		ID:           "",
		Email:        "admin@admin.com",
		Username:     "admin",
		PasswordHash: "admin",
	}
	storage := db.NewStorage(mongoDBClient, cfg.Mongodb.Collection, logger)
	user1Id, err := storage.Create(context.Background(), user1)
	if err != nil {
		panic(err)
	}
	logger.Info(user1Id)
	logger.Info("register user handler")
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server are listening unix socket :%s", socketPath)
	} else {
		logger.Info("listen tcp")
		logger.Infof("server are listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
