package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	//NewRequest returns a new incoming server Request to an http.Handler for testing.
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// ResponseRecorder is an implementation of http.ResponseWriter .Used here to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	res, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Check the response body is what we expect.
	if string(res) != "OK" {
		log.Fatal(err)
	}
}
