package form3

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	baseURLPath = "/api-v3"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "URL not found:")
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)
	client = CreateClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

func TestCreateClient(t *testing.T) {
	c := CreateClient(nil)

	if got, want := c.BaseURL.String(), defaultClientOptions.BaseEndpoint; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}

	c2 := CreateClient(nil)
	if c.httpClient == c2.httpClient {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}
}

func TestClient_createRequest(t *testing.T) {
	c := CreateClient(nil)

	inURL, outURL := "/foo", defaultClientOptions.BaseEndpoint+"/foo"
	inBody := MakeAccount("1", "2")
	req, _ := c.createRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("createRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	expectedBody := `{"data":{"id":"1","organisation_id":"2","type":"accounts","version":0}}`

	testBody(t, req, expectedBody)
	testHeader(t, req, "Content-Type", "application/json")

	// check that date header is set and in correct format
	_, err := time.Parse(time.RFC850, req.Header.Get("Date"))
	assert.Nil(t, err)
}

func TestClient_GET(t *testing.T) {
	c := CreateClient(nil)
	req, _ := c.createRequest("GET", "/foo", nil)
	testMethod(t, req, "GET")
}

func TestClient_POST(t *testing.T) {
	c := CreateClient(nil)
	req, _ := c.createRequest("POST", "/foo", nil)
	testMethod(t, req, "POST")
}

func TestClient_DELETE(t *testing.T) {
	c := CreateClient(nil)
	req, _ := c.createRequest("DELETE", "/foo", nil)
	testMethod(t, req, "DELETE")
}

func TestClient_Do(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/x", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{"data":{"A":"a"}}`)
	})

	req, _ := client.createRequest("GET", "/x", nil)
	body := new(foo)
	_, _ = client.Do(context.Background(), req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}
