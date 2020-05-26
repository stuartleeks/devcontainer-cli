package errors

var _ error = &StatusError{}

// StatusError is used to return informational errors
type StatusError struct {
	Message string
}

func (e *StatusError) Error() string { return e.Message }
