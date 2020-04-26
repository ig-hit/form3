package form3

import (
	"context"
	"fmt"
	"time"
)

const (
	accountsBaseEndpoint = "/organisation/accounts"
)

type AccountsService service

type Account struct {
	ID             string                `json:"id"`
	OrganisationID string                `json:"organisation_id"`
	Type           string                `json:"type"`
	Attributes     *AccountAttributes    `json:"attributes,omitempty"`
	Version        int                   `json:"version"`
	Relationships  *AccountRelationships `json:"relationships,omitempty"`
	CreatedOn      *time.Time            `json:"created_on,omitempty"`
	ModifiedOn     *time.Time            `json:"modified_on,omitempty"`
}

type ConfirmationOfPayeeAccount struct {
	*Account
	Attributes ConfirmationOfPayeeAttributes
}

type VirtualAccount struct {
	*Account
}

type AccountAttributes struct {
	Country                    string                      `json:"country"`
	BaseCurrency               string                      `json:"base_currency"`
	BankID                     string                      `json:"bank_id"`
	BankIDCode                 string                      `json:"bank_id_code"`
	AccountNumber              string                      `json:"account_number"`
	BIC                        string                      `json:"bic"`
	IBAN                       string                      `json:"iban"`
	CustomerID                 string                      `json:"customer_id"`
	Name                       string                      `json:"name"`
	PrivateIdentification      *PrivateIdentification      `json:"private_identification,omitempty"`
	OrganizationIdentification *OrganizationIdentification `json:"organization_identification,omitempty"`
}

type ConfirmationOfPayeeAttributes struct {
	*AccountAttributes

	AlternativeNames        []string `json:"alternative_names"`
	AccountClassification   string   `json:"account_classification"`
	JoinAccount             bool     `json:"join_account"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out"`
	SecondaryIdentification string   `json:"secondary_identification"`
	Switched                bool     `json:"switched"`
}

type PrivateIdentification struct {
	BirthDate      string   `json:"birth_date"`
	BirthCountry   string   `json:"birth_country"`
	Identification string   `json:"identification"`
	Address        []string `json:"address"`
	Country        string   `json:"country"`
	City           string   `json:"city"`
}

type OrganizationActor struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
	Residency string `json:"residency"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Country   string `json:"country"`
}

type OrganizationIdentification struct {
	Identification string            `json:"identification"`
	Actors         OrganizationActor `json:"actors"`
}

type MasterAccountRelation struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type AccountRelationships struct {
	MasterAccount []MasterAccountRelation
}

// Build new account
func MakeAccount(id, orgID string) *Account {
	return &Account{
		ID:             id,
		OrganisationID: orgID,
		Type:           "accounts",
	}
}

// Creates new account.
func (s *AccountsService) Create(ctx context.Context, data *Account) (*Account, *Response, error) {
	client := s.client

	req, err := client.POST(accountsBaseEndpoint, data)
	if err != nil {
		return nil, nil, nil
	}

	account := new(Account)
	resp, err := client.Do(ctx, req, account)
	if err != nil {
		return nil, resp, err
	}

	return account, resp, nil
}

// Retrieves account by ID.
func (s *AccountsService) ByID(ctx context.Context, id string) (*Account, *Response, error) {
	client := s.client
	url := fmt.Sprintf("%s/%s", accountsBaseEndpoint, id)

	req, err := client.GET(url, nil)
	if err != nil {
		return nil, nil, nil
	}

	account := new(Account)
	resp, err := client.Do(ctx, req, account)
	if err != nil {
		return nil, resp, err
	}

	return account, resp, nil
}

// Specify pagination options.
type AccountListOptions struct {
	// Page number being requested: int or first|last
	Number string

	// Size of the page being requested
	Size int
}

// Get list of accounts.
func (s *AccountsService) List(ctx context.Context, options *AccountListOptions) ([]*Account, *Response, error) {
	client := s.client

	reqUrl := accountsBaseEndpoint
	if options != nil {
		reqUrl = fmt.Sprintf("%s?page[number]=%s&page[size]=%d", accountsBaseEndpoint, options.Number, options.Size)
	}

	req, err := client.GET(reqUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	accounts := new([]*Account)
	resp, err := client.Do(ctx, req, accounts)
	if err != nil {
		return nil, resp, err
	}

	return *accounts, resp, nil
}

// Delete account by id and version
func (s *AccountsService) Delete(ctx context.Context, id string, version int) (*Response, error) {
	client := s.client
	url := fmt.Sprintf("%s/%s?version=%d", accountsBaseEndpoint, id, version)

	req, err := client.DELETE(url, nil)
	if err != nil {
		return nil, nil
	}

	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
