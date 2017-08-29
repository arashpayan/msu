package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

// LogCalls ...
var LogCalls = false

// Handler ...
func Handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if LogCalls {
			log.Printf("%s %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)
		}

		defer func() {
			if r := recover(); r != nil {
				s := make([]byte, 2048)
				numBytes := runtime.Stack(s, false)
				stack := s[:numBytes]
				err := fmt.Errorf("recovered - %v\n%s", r, string(stack))
			}
		}()
		next.ServeHTTP(w, r)
	}
}

// SendResponse ...
func SendResponse(w http.ResponseWriter, response interface{}, httpCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)

	enc := json.NewEncoder(w)
	err := enc.Encode(response)
	if err != nil {
		panic(err)
	}
}

// SendErr ...
func SendErr(w http.ResponseWriter, msg string, httpCode int, apiCode int) {
	SendResponse(
		w,
		struct {
			Msg  string `json:"error_message"`
			Code int    `json:"error_code"`
		}{Msg: msg, Code: apiCode},
		httpCode)
}

// SendInternalErr ...
func SendInternalErr(w http.ResponseWriter, err error) {
	SendErr(w, "Internal server error", http.StatusInternalServerError, ErrorInternal)

	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}
		file = filepath.Base(file)
		log.Printf("%s:%d %v", file, line, err)
	}
}

// SendBadReqCode ...
func SendBadReqCode(w http.ResponseWriter, msg string, apiCode int) {
	SendErr(w, msg, http.StatusBadRequest, apiCode)
}

// SendBadReq ...
func SendBadReq(w http.ResponseWriter, msg string) {
	SendBadReqCode(w, msg, ErrorBadRequest)
}

// SendSuccess ...
func SendSuccess(w http.ResponseWriter, resp interface{}) {
	if resp == nil {
		resp = struct{}{}
	}

	SendResponse(w, resp, http.StatusOK)
}
