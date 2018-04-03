package msu

import (
	"bytes"
	"log"
)

// LogFlags ...
func LogFlags() int {
	return log.Ldate | log.Ltime | log.Lshortfile
}

// TLSHandshakeFilter can be used as the ErrorLog value in an HTTP server.
// All it does is filter out TLS handshake errors from the log and then
// pass the rest on to the log
type TLSHandshakeFilter struct{}

func (dl *TLSHandshakeFilter) Write(p []byte) (int, error) {
	if bytes.Contains(p, []byte("TLS handshake error from")) {
		return len(p), nil // lie to the caller
	}

	log.Printf("%s", p)
	return len(p), nil
}
