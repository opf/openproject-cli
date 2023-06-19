package requests

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

var client *http.Client
var host *url.URL

func Init(hostUrl string) {
	client = &http.Client{}
	
	var err error
	host, err = url.Parse(hostUrl)
	if err != nil {
		log.Fatalf("Invalid host url '%s'.", hostUrl)
	}
}

func Get(path string) []byte {
	if client == nil {
		log.Fatal("Cannot execute requests without initializing request client first.")
	}
	
	host.Path = path
	resp, err := client.Get(host.String())
	if err != nil {
		log.Fatalf("Error execting request: %+v", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	return body
}
