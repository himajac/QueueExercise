package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Logger Instance
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.WithPrefix(logger, "date", log.DefaultTimestampUTC)

	//q := NewQueue()      //Slice based Queue
	linkedListQ := NewLinkedListQueue(logger)

	//Create handler Instance
	h := newHandler(linkedListQ, logger)

	//Create Router Instance
	router := newRouter(&h)

	//Create http Server
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	//To store fatalErrors while server is running
	fatalErrorChan := make(chan error)
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	logger.Log("level", "info", "msg", "starting server", "port", "8080")

	//start Server
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fatalErrorChan <- errors.Wrap(err, "server failed")
		}
	}()

	//blocking
	select {
	case sig := <-sigChan:
		logger.Log("level", "info", "msg", "received os signal", "signal", sig.String())
		break
	case err := <-fatalErrorChan:
		logger.Log("level", "error", "msg", "received fatal error, shutting down", "error", err.Error())
		os.Exit(1)
	}

	logger.Log("level", "info", "msg", "waiting on open connections to finish")

	err := server.Shutdown(context.Background())
	if err != nil {
		logger.Log("level", "error", "msg", "failed to shutdown server")
		os.Exit(1)
	}
	logger.Log("level", "info", "msg", "shutdown successful.")
	os.Exit(0)
}
