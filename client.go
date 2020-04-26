package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClientOptions struct {
	// milliseconds
	Timeout int

	// http://localhost:8080/v1
	BaseEndpoint string
}

// A Client manages communication with the API.
type Client struct {
	httpClient *http.Client

	BaseURL *url.URL
}

type service struct {
	client *Client
}

type body struct {
	// wraps payload into data object
	Data interface{} `json:"data"`
}

type bodyError struct {
	ErrorMessage string `json:"error_message"`
}

func (c *Client) GET(url string, body interface{}) (*http.Request, error) {
	return c.createRequest("GET", url, body)
}

func (c *Client) POST(url string, body interface{}) (*http.Request, error) {
	return c.createRequest("POST", url, body)
}

func (c *Client) DELETE(url string, body interface{}) (*http.Request, error) {
	return c.createRequest("DELETE", url, body)
}

func (c *Client) createRequest(method, url string, payload interface{}) (*http.Request, error) {
	reqUrl := fmt.Sprintf("%s/%s", strings.TrimRight(c.BaseURL.String(), "/"), strings.TrimLeft(url, "/"))

	var data io.Reader
	if payload != nil {
		body := body{Data: payload}
		dataJson, _ := json.Marshal(body)
		data = bytes.NewBuffer(dataJson)
	}

	req, err := http.NewRequest(method, reqUrl, data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Date", time.Now().Format(time.RFC850))

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, target interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("nil context is not allowed")
	}

	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	respText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	// check if there're any errors
	be := &bodyError{}
	_ = json.Unmarshal(respText, be)
	if be.ErrorMessage != "" {
		return response, errors.New(be.ErrorMessage)
	}

	// unpack into target
	if target != nil {
		body := &body{}
		err = json.Unmarshal(respText, body)
		if err != nil {
			return response, err
		}

		encoded, err := json.Marshal(body.Data)
		if err == nil {
			err = json.Unmarshal(encoded, target)
		}
	}

	return response, err
}

type Response struct {
	*http.Response
}

func newResponse(r *http.Response) *Response {
	return &Response{Response: r}
}
