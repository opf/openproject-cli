package requests

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/opf/openproject-cli/components/printer"
)

var client *http.Client
var host *url.URL
var token string

type RequestBody struct {
	ContentType string
	Body        io.Reader
}

func Init(hostUrl *url.URL, tokenValue string) {
	client = &http.Client{}
	host = hostUrl
	token = tokenValue
}

func Get(path string, query *Query) (code int, body []byte) {
	return Do("GET", path, query, nil)
}

func Do(method string, path string, query *Query, reqBody *RequestBody) (status int, response []byte) {
	if client == nil {
		printer.ErrorText("Cannot execute requests without initializing request client first. Run `op login`")
	}

	requestUrl := *host
	requestUrl.Path = path
	if query != nil {
		requestUrl.RawQuery = query.String()
	}

	var body io.Reader
	if reqBody != nil {
		body = reqBody.Body
	}

	request, err := http.NewRequest(
		strings.ToUpper(method),
		requestUrl.String(),
		body,
	)
	if err != nil {
		printer.Error(err)
	}

	if reqBody != nil {
		request.Header.Add("Content-Type", reqBody.ContentType)
	}

	if len(token) > 0 {
		request.SetBasicAuth("apikey", token)
	}

	resp, err := client.Do(request)
	if err != nil {
		printer.Error(err)
	}

	defer func() { _ = resp.Body.Close() }()

	response, err = io.ReadAll(resp.Body)
	if err != nil {
		printer.Error(err)
	}

	return resp.StatusCode, response
}

func IsSuccess(code int) bool {
	return code >= 200 && code <= 299
}
