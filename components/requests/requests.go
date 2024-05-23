package requests

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
)

var client *http.Client
var host *url.URL
var token string

type RequestData struct {
	ContentType string
	Body        io.Reader
}

func Init(hostUrl *url.URL, tokenValue string) {
	client = &http.Client{}
	host = hostUrl
	token = tokenValue
}

func Get(path string, query *Query) (responseBody []byte, err error) {
	loadingFunc := func() ([]byte, error) { return Do("GET", path, query, nil) }

	return printer.WithSpinner(loadingFunc)
}

func Post(path string, requestData *RequestData) (responseBody []byte, err error) {
	loadingFunc := func() ([]byte, error) { return Do("POST", path, nil, requestData) }

	return printer.WithSpinner(loadingFunc)
}

func Patch(path string, requestBody *RequestData) (responseBody []byte, err error) {
	loadingFunc := func() ([]byte, error) { return Do("PATCH", path, nil, requestBody) }

	return printer.WithSpinner(loadingFunc)
}

func Do(method string, path string, query *Query, requestData *RequestData) (responseBody []byte, err error) {
	if client == nil || hostUnitialised() {
		return nil, errors.Custom("Cannot execute requests without initializing request client first. Run `op login`")
	}

	requestUrl := *host
	requestUrl.Path += path
	if query != nil {
		requestUrl.RawQuery = query.String()
	}

	var body io.Reader
	if requestData != nil {
		body = requestData.Body
	}

	request, err := http.NewRequest(
		strings.ToUpper(method),
		requestUrl.String(),
		body,
	)
	if err != nil {
		return nil, err
	}

	if requestData != nil {
		request.Header.Add("Content-Type", requestData.ContentType)
	}

	if len(token) > 0 {
		request.SetBasicAuth("apikey", token)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !isSuccess(resp.StatusCode) {
		return nil, errors.NewResponseError(resp.StatusCode, response)
	}

	return response, nil
}

func isSuccess(code int) bool {
	return code >= 200 && code <= 299
}

func hostUnitialised() bool {
	return len(host.Scheme) == 0 || len(host.Host) == 0
}
