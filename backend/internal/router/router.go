package router

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"neatly/internal/context"
	"neatly/pkg/logging"
	"neatly/pkg/metric/handler"
	"neatly/pkg/shutdown"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

func Init() {
	logger := logging.GetLogger()
	logger.Println("Initializing application router")

	router := httprouter.New()

	router.GET(handler.HEARTBEAT_URL, handler.Heartbeat)

	ctx := context.GetInstance()
	cfg := ctx.Config

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
