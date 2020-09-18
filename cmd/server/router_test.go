package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func do(h handler, method string, url string, header http.Header, body interface{}) (*httptest.ResponseRecorder, *http.Request) {
	testRouter := newRouter(&h)
	b, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	r := httptest.NewRequest(method, url, bytes.NewBuffer(b))
	r.Header = header

	wr := httptest.NewRecorder()
	testRouter.ServeHTTP(wr, r)
	return wr, r
}
