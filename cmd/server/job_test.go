package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestJobListQueue_Enqueue2(t *testing.T) {
	j := job{Type: "TIME_CRITICAL"}
	body, err := json.Marshal(j)

	req, err := http.NewRequest("POST", "/jobs/enqueue", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	q := NewLinkedListQueue(log.NewNopLogger())
	h := newHandler(q, log.NewNopLogger())

	// ResponseRecorder is an implementation of http.ResponseWriter .Used here to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.enqueue)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	res, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Error(err.Error())
	}
	var enqueueResponse jobIdResponse
	err = json.Unmarshal(res, &enqueueResponse)
	if err != nil {
		t.Error(err.Error())
	}
	if enqueueResponse.JobId == 0 {
		t.Errorf("jobId returned from Enqueue cannot be 0: got %v", enqueueResponse.JobId)
	}
}

func TestJobListQueue_Dequeue2(t *testing.T) {
	req, err := http.NewRequest("GET", "/jobs/dequeue", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := NewLinkedListQueue(log.NewNopLogger())
	item1 := job{Type: "TIME_CRITICAL"}
	id1, _ := q.Enqueue(&item1)

	h := newHandler(q, log.NewNopLogger())

	// ResponseRecorder is an implementation of http.ResponseWriter .Used here to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.dequeue)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	res, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Error(err.Error())
	}

	var dequeueResponse job
	err = json.Unmarshal(res, &dequeueResponse)
	if err != nil {
		t.Error(err.Error())
	}

	if dequeueResponse.Id != id1 {
		t.Errorf("jobId returned from Enqueue cannot be 0: got %v", dequeueResponse.Id)
	}

	//The next call to dequeue should return 500. Since no jobs are  available to dequeue
	rrEmpty := httptest.NewRecorder()
	handler = http.HandlerFunc(h.dequeue)
	handler.ServeHTTP(rrEmpty, req)

	// Check the status code is what we expect.
	if status := rrEmpty.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestJobListQueue_Conclude2(t *testing.T) {
	q := NewLinkedListQueue(log.NewNopLogger())
	item1 := job{Type: "TIME_CRITICAL"}
	id1, _ := q.Enqueue(&item1)

	url := fmt.Sprintf("/jobs/%s/conclude", strconv.Itoa(id1))
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	//Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"job_id": strconv.Itoa(id1),
	}
	req = mux.SetURLVars(req, vars)

	h := newHandler(q, log.NewNopLogger())

	// ResponseRecorder is an implementation of http.ResponseWriter .Used here to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.conclude)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	res, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Error(err.Error())
	}

	var concludeResponse jobIdResponse
	err = json.Unmarshal(res, &concludeResponse)
	if err != nil {
		t.Error(err.Error())
	}

	if concludeResponse.JobId != id1 {
		t.Errorf("jobId returned from Enqueue cannot be 0: got %v", concludeResponse.JobId)
	}
}

func TestJobListQueue_GetJob2(t *testing.T) {
	q := NewLinkedListQueue(log.NewNopLogger())
	item1 := job{Type: "TIME_CRITICAL"}
	id1, _ := q.Enqueue(&item1)

	url := fmt.Sprintf("/jobs/%s", strconv.Itoa(id1))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	//Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"job_id": strconv.Itoa(id1),
	}
	req = mux.SetURLVars(req, vars)

	h := newHandler(q, log.NewNopLogger())

	// ResponseRecorder is an implementation of http.ResponseWriter .Used here to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.getJob)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	res, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Error(err.Error())
	}

	var jobResponse job
	err = json.Unmarshal(res, &jobResponse)
	if err != nil {
		t.Error(err.Error())
	}

	if jobResponse.Id != id1 {
		t.Errorf("jobId returned from Enqueue cannot be 0: got %v", jobResponse.Id)
	}
}
