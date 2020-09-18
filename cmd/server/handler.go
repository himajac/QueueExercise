package main

import (
	"github.com/go-kit/kit/log"
)

type handler struct {
	queue  Queue
	logger log.Logger
}

func newHandler(queue Queue, log log.Logger) handler {
	return handler{queue, log}
}
