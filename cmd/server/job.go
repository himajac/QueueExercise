package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
)

//Each item in the queue is of type job
type job struct {
	Id     int    `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	mutex  sync.Mutex
}

type jobIdResponse struct {
	JobId int `json:"jobId"`
}

func (h *handler) enqueue(w http.ResponseWriter, r *http.Request) {
	var req job
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Log("level", "error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jobId, err := h.queue.Enqueue(&req)
	if err != nil {
		h.logger.Log("level", "error", "msg", "Not able to queue job", "error", err.Error())
		Respond(w, http.StatusInternalServerError, err.Error())
	}

	enqueueResponse := jobIdResponse{JobId: jobId}
	Respond(w, http.StatusCreated, enqueueResponse)
	return
}

func (h *handler) dequeue(w http.ResponseWriter, r *http.Request) {
	jobDeque, err := h.queue.Dequeue()
	if err != nil {
		h.logger.Log("level", "error", "error", err.Error())
		Respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	Respond(w, http.StatusOK, jobDeque)
	return
}

func (h *handler) conclude(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["job_id"]
	if !ok {
		h.logger.Log("level", "error", "msg", "could not get JobId from the request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jobID, _ := strconv.Atoi(id)

	err := h.queue.Conclude(jobID)
	if err != nil {
		Respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	concludeResponse := jobIdResponse{JobId: jobID}
	Respond(w, http.StatusOK, concludeResponse)
	return
}

func (h *handler) getJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["job_id"]
	if !ok {
		h.logger.Log("level", "error", "msg", "could not get JobId from the request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jobID, _ := strconv.Atoi(id)
	jobDetails, err := h.queue.GetJob(jobID)
	if err != nil {
		Respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	Respond(w, http.StatusOK, jobDetails)
	return
}

func (h *handler) getJobs(w http.ResponseWriter, r *http.Request) {
	jobDetails, err := h.queue.GetJobs()
	if err != nil {
		Respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	Respond(w, http.StatusOK, jobDetails)
	return
}

func (h *handler) remove(w http.ResponseWriter, r *http.Request) {
	jobId, err := h.queue.Remove()
	if err != nil {
		Respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	deletedResponse := jobIdResponse{JobId: jobId}
	Respond(w, http.StatusOK, deletedResponse)
	return
}
