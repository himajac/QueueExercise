package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter(h *handler) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	//Create a subRouter for all the paths with prefix jobs
	jobsRouter := router.PathPrefix("/jobs").Subrouter()
	jobsRouter.HandleFunc("/enqueue", h.enqueue).Methods(http.MethodPost)
	jobsRouter.HandleFunc("/dequeue", h.dequeue).Methods(http.MethodGet)
	jobsRouter.HandleFunc("/{job_id}/conclude", h.conclude).Methods(http.MethodPost)
	jobsRouter.HandleFunc("/{job_id}", h.getJob).Methods(http.MethodGet)
	jobsRouter.HandleFunc("", h.remove).Methods(http.MethodDelete)
	jobsRouter.HandleFunc("", h.getJobs).Methods(http.MethodGet)
	return router
}
