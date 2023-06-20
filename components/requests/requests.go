package requests

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

var client *http.Client
var host *url.URL
var token string

func Init(hostUrl, tokenValue string) {
	client = &http.Client{}

	var err error
	host, err = url.Parse(hostUrl)
	if err != nil {
		log.Fatalf("Invalid host url '%s'.", hostUrl)
	}

	token = tokenValue
}

func Get(path string) (success bool, response []byte) {
	if client == nil {
		log.Fatal("Cannot execute requests without initializing request client first. Run `op login`")
	}

	requestUrl := *host
	requestUrl.Path = path
	request, _ := http.NewRequest("GET", requestUrl.String(), nil)
	if len(token) > 0 {
		request.SetBasicAuth("apikey", token)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error execting request: %+v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	return isSuccess(resp), body
}

func isSuccess(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode <= 299
}
