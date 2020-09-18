package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func Respond(w http.ResponseWriter, statusCode int, payload interface{}) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload")
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(statusCode)
	w.Write(bytes)

	return nil
}

func RespondOK(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusOK, payload)
}
