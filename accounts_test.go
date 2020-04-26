package form3

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func setupAccounts() (service *AccountsService, mux *http.ServeMux, serverURL string, teardown func()) {
	client, mux, serverURL, teardown := setup()
	service = CreateAccountsService(client)
	return service, mux, serverURL, teardown
}

func TestMakeAccount(t *testing.T) {
	acc := MakeAccount("1", "2")
	assert.Equal(t, "accounts", acc.Type)
}

func TestAccountsService_Create(t *testing.T) {
	service, mux, _, teardown := setupAccounts()
	defer teardown()

	mux.HandleFunc("/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(
			w,
			`{"data":{"id":"1","organisation_id":"2","type":"accounts","version":0}}`,
		)
	})

	account := MakeAccount("1", "2")

	saved, resp, err := service.Create(context.Background(), account)
	assert.Nil(t, err)

	expected := MakeAccount("1", "2")
	assert.Equal(t, expected, saved)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestAccountsService_ByID(t *testing.T) {
	service, mux, _, teardown := setupAccounts()
	defer teardown()

	mux.HandleFunc("/organisation/accounts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(
			w,
			`{"data":{"id":"1","organisation_id":"2","type":"accounts","version":0}}`,
		)
	})

	saved, resp, err := service.ByID(context.Background(), "1")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "1", saved.ID)
}

func TestAccountsService_List(t *testing.T) {
	service, mux, _, teardown := setupAccounts()
	defer teardown()

	mux.HandleFunc("/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		items := []string{
			`{"id":"x-1","organisation_id":"2","type":"accounts","version":0}`,
			`{"id":"x-2","organisation_id":"2","type":"accounts","version":0}`,
			`{"id":"x-3","organisation_id":"2","type":"accounts","version":0}`,
		}
		_, _ = fmt.Fprint(w, fmt.Sprintf(`{"data":[%s]}`, strings.Join(items, ",")))
	})

	list, resp, err := service.List(context.Background(), nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, 3, len(list))
	for i, acc := range list {
		assert.Equal(t, fmt.Sprintf("x-%d", i+1), acc.ID)
	}
}

func TestAccountsService_ListWithOptions(t *testing.T) {
	service, mux, _, teardown := setupAccounts()
	defer teardown()

	mux.HandleFunc("/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		query := r.URL.Query()
		pageNumber, _ := strconv.Atoi(query.Get("page[number]"))
		pageSize, _ := strconv.Atoi(query.Get("page[size]"))
		offset := (pageNumber - 1) * pageSize

		all := []string{
			`{"id":"x-1","organisation_id":"2","type":"accounts","version":0}`,
			`{"id":"x-2","organisation_id":"2","type":"accounts","version":0}`,
			`{"id":"x-3","organisation_id":"2","type":"accounts","version":0}`,
			`{"id":"x-4","organisation_id":"2","type":"accounts","version":0}`,
		}

		items := all[offset : offset+pageSize]

		_, _ = fmt.Fprint(w, fmt.Sprintf(`{"data":[%s]}`, strings.Join(items, ",")))
	})

	options := &AccountListOptions{
		Number: "3",
		Size:   1,
	}
	list, resp, err := service.List(context.Background(), options)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "x-3", list[0].ID)
}

func TestAccountsService_Delete(t *testing.T) {
	service, mux, _, teardown := setupAccounts()
	defer teardown()

	mux.HandleFunc("/organisation/accounts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := service.Delete(context.Background(), "1", 0)
	assert.Nil(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}
