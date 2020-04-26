// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/ig-hit/form3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// Checks that account is successfully created,
// has correct status code and response.
func TestAccountsService_Create(t *testing.T) {
	service := form3.CreateAccountsService(nil)
	account := createAccount(t)

	expected := form3.MakeAccount(account.ID, account.OrganisationID)
	expected.Attributes = &form3.AccountAttributes{
		Country:                    account.Attributes.Country,
		BaseCurrency:               account.Attributes.BaseCurrency,
		BankID:                     account.Attributes.BankID,
		BankIDCode:                 account.Attributes.BankIDCode,
		CustomerID:                 account.Attributes.CustomerID,
	}

	ctx := context.Background()
	saved, resp, err := service.Create(ctx, account)
	assert.Nil(t, err)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, saved.CreatedOn)

	expected.CreatedOn = saved.CreatedOn
	expected.ModifiedOn = saved.ModifiedOn

	assert.NotNil(t, saved)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, expected, saved)
}

func TestAccountsService_CreateError(t *testing.T) {
	service := form3.CreateAccountsService(nil)
	account := createAccount(t)
	account.ID = "x"

	ctx := context.Background()
	saved, resp, err := service.Create(ctx, account)

	assert.Nil(t, saved)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, err.Error(), `id in body must be of type uuid: "x"`)
}

func TestAccountsService_CreateDataError(t *testing.T) {
	service := form3.CreateAccountsService(nil)
	account := createAccount(t)
	account.Attributes.Country = "XXX"

	ctx := context.Background()
	saved, resp, err := service.Create(ctx, account)

	assert.Nil(t, saved)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, err.Error(), `country in body should match '^[A-Z]{2}$'`)
}

// Checks that account can be retrieved by id,
// has correct http status code and response.
func TestAccountsService_Get(t *testing.T) {
	service := form3.CreateAccountsService(nil)
	account := createAccount(t)
	ctx := context.Background()
	_, _, err := service.Create(ctx, account)
	assert.Nil(t, err)

	actual, resp, err := service.ByID(ctx, account.ID)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expected := form3.MakeAccount(account.ID, account.OrganisationID)
	expected.Attributes = &form3.AccountAttributes{
		Country:                    account.Attributes.Country,
		BaseCurrency:               account.Attributes.BaseCurrency,
		BankID:                     account.Attributes.BankID,
		BankIDCode:                 account.Attributes.BankIDCode,
		CustomerID:                 account.Attributes.CustomerID,
	}
	expected.CreatedOn = actual.CreatedOn
	expected.ModifiedOn = actual.ModifiedOn

	assert.Equal(t, expected, actual)
}

func TestAccountsService_GetNotFound(t *testing.T) {
	service := form3.CreateAccountsService(nil)
	nonExistentID := form3.CreateUUID()
	actual, resp, err := service.ByID(context.Background(), nonExistentID)
	assert.NotNil(t, err)
	assert.Nil(t, actual)
	assert.Contains(t, err.Error(), fmt.Sprintf("record %s does not exist", nonExistentID))
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// Checks that account can be successfully deleted
func TestAccountsService_Delete(t *testing.T) {
	service := form3.CreateAccountsService(nil)
	account := createAccount(t)

	ctx := context.Background()
	_, _, err := service.Create(ctx, account)
	assert.Nil(t, err)

	resp, err := service.Delete(ctx, account.ID, account.Version)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// Checks that non-existent ID deletion returns corresponding http status code.
func TestAccountsService_DeleteNonExistent(t *testing.T) {
	service := form3.CreateAccountsService(nil)
	account := createAccount(t)

	ctx := context.Background()
	_, _, err := service.Create(ctx, account)
	assert.Nil(t, err)

	resp, err := service.Delete(ctx, form3.CreateUUID(), account.Version)

	assert.Nil(t, err)
	// todo: #clarify
	// according to https://api-docs.form3.tech/api.html#organisation-accounts-delete code should be 404
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// Checks that ID must be in UUID format.
func TestAccountsService_DeleteInvalid(t *testing.T) {
	service := form3.CreateAccountsService(nil)
	account := createAccount(t)

	ctx := context.Background()
	_, _, err := service.Create(ctx, account)
	assert.Nil(t, err)

	invalidID := "1"
	resp, err := service.Delete(ctx, invalidID, account.Version)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// Checks that wrong version deletion responds with the correct http status code
func TestAccountsService_DeleteWrongVersion(t *testing.T) {
	service := form3.CreateAccountsService(nil)
	account := createAccount(t)

	ctx := context.Background()
	_, _, err := service.Create(ctx, account)
	assert.Nil(t, err)

	resp, err := service.Delete(ctx, account.ID, 1)
	assert.NotNil(t, err)
	// todo: #clarify
	// according to https://api-docs.form3.tech/api.html#organisation-accounts-delete code should be 409
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Contains(t, err.Error(), "invalid version")
}

// Checks that correct records are returned without pagination options
func TestAccountsService_List(t *testing.T) {
	service := form3.CreateAccountsService(nil)
	qty := 10
	populated := populate(t, qty)

	list, resp, err := service.List(context.Background(), nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	returnedIDs := make(map[string]*form3.Account)
	for _, a := range list {
		returnedIDs[a.ID] = a
	}

	for id, p := range populated {
		saved, found := returnedIDs[id]
		assert.True(t, found, fmt.Sprintf("Populated ID %s not found in returned data", p.ID))
		assert.Equal(t, p, saved)
	}
}

// Checks that correct number of records is returned with pagination options
func TestAccountsService_ListWithPagination(t *testing.T) {
	service := form3.CreateAccountsService(nil)
	pageSize := 12
	// populate a little more
	_ = populate(t, 3 * pageSize)

	options := &form3.AccountListOptions{
		Number: "first",
		Size:   pageSize,
	}

	list, resp, err := service.List(context.Background(), options)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, pageSize, len(list))
}

// Checks that list can return 0-length accounts
func TestAccountsService_ListWithPaginationOutOfBounds(t *testing.T) {
	service := form3.CreateAccountsService(nil)

	options := &form3.AccountListOptions{
		Number: "100",
		Size:   100,
	}

	list, resp, err := service.List(context.Background(), options)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 0, len(list))
}

// Helper function.
func createAccount(t *testing.T) *form3.Account {
	account := form3.MakeAccount(form3.CreateUUID(), form3.CreateUUID())
	account.Attributes = &form3.AccountAttributes{
		Country:                    "DE",
		BaseCurrency:               "EUR",
		BankID:                     "12345678",
		BankIDCode:                 "DEBLZ",
		CustomerID:                 "XXX-3",
		Name:                       "Jose2 Sanchez",
	}

	return account
}

func populate(t *testing.T, n int) map[string]*form3.Account {
	service := form3.CreateAccountsService(nil)
	ctx := context.Background()
	as := make(map[string]*form3.Account, 0)

	for i := 0; i < n; i++ {
		a := createAccount(t)
		saved, _, _ := service.Create(ctx, a)
		as[saved.ID] = saved
	}

	return as
}
