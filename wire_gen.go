package form3

import (
	"net/http"
	"net/url"
	"time"
)

var defaultClientOptions = &ClientOptions{
	Timeout:      3000,
	BaseEndpoint: "http://localhost:8080/v1",
}

func provideHTTPClient(timeout int) *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
}

func CreateClient(options *ClientOptions) *Client {
	if options == nil {
		options = defaultClientOptions
	}

	baseURL, _ := url.Parse(options.BaseEndpoint)

	return &Client{
		httpClient: provideHTTPClient(options.Timeout),
		BaseURL:    baseURL,
	}
}

func CreateAccountsServiceWithOptions(options *ClientOptions) *AccountsService {
	return &AccountsService{
		client: CreateClient(options),
	}
}

func CreateAccountsService(client *Client) *AccountsService {
	if client == nil {
		client = CreateClient(nil)
	}
	return &AccountsService{
		client: client,
	}
}
