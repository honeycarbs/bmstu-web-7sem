package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"neatly/internal/handlers/auth"
	"neatly/internal/session"
	"neatly/pkg/cache/freecache"
	"neatly/pkg/handlers/metric"
	"neatly/pkg/logging"
	"neatly/pkg/shutdown"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

func main() {
	logging.Init()

	logger := logging.GetLogger()
	logger.Println("Logger initialized.")

	logger.Println("Application app-context initialized.")

	router := httprouter.New()
	cfg := session.GetConfig()

	refreshTokenCache := freecache.NewCacheRepo(104857600)

	authHandler := auth.Handler{RTCache: refreshTokenCache, Logger: logger}
	authHandler.Register(router)

	metricHandler := metric.Handler{Logger: logger}
	metricHandler.Register(router)

	logger.Println("Application started.")

	start(router, logger, cfg)
}

func start(router *httprouter.Router, logger logging.Logger, cfg *session.Config) {
	var server *http.Server
	var listener net.Listener

	logger.Infof(
		"Bind application to : %s:%s",
		cfg.Listen.BindIP,
		cfg.Listen.Port,
	)

	var err error
	listener, err = net.Listen("tcp", fmt.Sprintf(
		"%s:%s",
		cfg.Listen.BindIP,
		cfg.Listen.Port,
	))
	if err != nil {
		logger.Fatal(err)
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server)

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
