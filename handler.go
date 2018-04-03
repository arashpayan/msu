package msu

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

// LogHTTPCalls ...
var LogHTTPCalls = false

// Handler ...
func Handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if LogHTTPCalls {
			log.Printf("%s %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)
		}

		defer func() {
			if r := recover(); r != nil {
				s := make([]byte, 2048)
				numBytes := runtime.Stack(s, false)
				stack := s[:numBytes]
				err := fmt.Errorf("recovered - %v\n%s", r, string(stack))
				InternalErr(w, err)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

// Respond ...
func Respond(w http.ResponseWriter, response interface{}, httpCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)

	enc := json.NewEncoder(w)
	err := enc.Encode(response)
	if err != nil {
		panic(err)
	}
}

// Err ...
func Err(w http.ResponseWriter, msg string, httpCode int, apiCode int) {
	Respond(
		w,
		struct {
			Msg  string `json:"error_message"`
			Code int    `json:"error_code"`
		}{Msg: msg, Code: apiCode},
		httpCode)
}

// InternalErr ...
func InternalErr(w http.ResponseWriter, err error) {
	Err(w, "Internal server error", http.StatusInternalServerError, ErrorInternal)

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

// BadReqCode ...
func BadReqCode(w http.ResponseWriter, msg string, apiCode int) {
	Err(w, msg, http.StatusBadRequest, apiCode)
}

// BadReq ...
func BadReq(w http.ResponseWriter, msg string) {
	BadReqCode(w, msg, ErrorBadRequest)
}

// Success ...
func Success(w http.ResponseWriter, resp interface{}) {
	if resp == nil {
		resp = struct{}{}
	}

	Respond(w, resp, http.StatusOK)
}
