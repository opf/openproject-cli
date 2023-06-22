package requests

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/opf/openproject-cli/components/printer"
)

var client *http.Client
var host *url.URL
var token string

func Init(hostUrl, tokenValue string) {
	client = &http.Client{}

	var err error
	host, err = url.Parse(hostUrl)
	if err != nil {
		printer.Error(err)
	}

	token = tokenValue
}

func Get(path string, query *Query) (code int, body []byte) {
	return Do("GET", path, query, nil)
}

func Do(method string, path string, query *Query, reqBody []byte) (code int, body []byte) {
	if client == nil {
		printer.ErrorText("Cannot execute requests without initializing request client first. Run `op login`")
	}

	requestUrl := *host
	requestUrl.Path = path
	if query != nil {
		requestUrl.RawQuery = query.String()
	}

	request, err := http.NewRequest(
		strings.ToUpper(method),
		requestUrl.String(),
		bytes.NewReader(reqBody),
	)
	if err != nil {
		printer.Error(err)
	}

	if reqBody != nil {
		request.Header.Add("Content-Type", "application/json")
	}

	if len(token) > 0 {
		request.SetBasicAuth("apikey", token)
	}

	resp, err := client.Do(request)
	if err != nil {
		printer.Error(err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		printer.Error(err)
	}

	return resp.StatusCode, body
}

func IsSuccess(code int) bool {
	return code >= 200 && code <= 299
}
