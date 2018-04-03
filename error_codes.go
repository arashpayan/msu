package msu

// Error codes used in http responses. Values less than 1000 are reserved
// for msu. Custom codes should start at 1000.
const (
	ErrorNone       int = 0
	ErrorInternal       = 1
	ErrorBadRequest     = 2
)
