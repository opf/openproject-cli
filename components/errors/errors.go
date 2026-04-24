package errors

import stderrors "errors"

type customError struct {
	text string
}

func (err *customError) Error() string {
	return err.text
}

func Custom(text string) error {
	return &customError{text: text}
}

func IsCustom(err error) bool {
	var target *customError
	return stderrors.As(err, &target)
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
