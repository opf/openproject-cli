package errors

import "errors"

type OpError int

const (
	NotInitialize OpError = iota
	HttpExecutionError
	HttpResponse
)

func New(err OpError) error {
	return newError(err, nil)
}

func NewWithReason(err OpError, reason error) error {
	return newError(err, reason)
}

func newError(err OpError, reason error) error {
	switch err {
	case NotInitialize:
		return &NotInitialized{}
	default:
		if reason == nil {
			return errors.New("unknown error occured")
		} else {
			return reason
		}
	}
}

type NotInitialized struct{}

func (err *NotInitialized) Error() string {
	return "Cannot execute requests without initializing request client first. Run `op login`."
}

type HttpResponseError struct {
	Status int
	Body   []byte
}

func (err *HttpResponseError) Error() string {
	return "http error"
}
