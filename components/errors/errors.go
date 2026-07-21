package errors

import "errors"

// ErrHandled signals that an error was already reported to the user via the
// printer. Commands return it from RunE so the process exits nonzero without
// printing the error a second time in main.
var ErrHandled = errors.New("error already reported")

type customError struct {
	text string
}

func (err *customError) Error() string {
	return err.text
}

func Custom(text string) error {
	return &customError{text: text}
}

type ResponseError struct {
	status   int
	response []byte
}

func (err *ResponseError) Error() string {
	return string(err.response)
}

func (err *ResponseError) Status() int {
	return err.status
}

func (err *ResponseError) Response() []byte {
	return err.response
}

func NewResponseError(status int, response []byte) error {
	return &ResponseError{status: status, response: response}
}
