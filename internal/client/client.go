package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/vestamart/homework/internal/domain"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	url        string
	token      string
}

type request struct {
	Token string `json:"token"`
	SKU   int64  `json:"sku"`
}

func NewClient(client *http.Client, url, token string) *Client {
	return &Client{
		httpClient: client,
		url:        url,
		token:      token,
	}
}
func (c *Client) ExistItem(sku int64) (bool, error) {
	jsonBody, err := json.Marshal(request{Token: c.token, SKU: sku})
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, domain.ErrSkuNotExist
	}
	return true, nil
}

func (c *Client) GetProductHandler(sku int64) (*domain.ClientResponse, error) {
	jsonBody, err := json.Marshal(request{Token: c.token, SKU: sku})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, domain.ErrSkuNotExist
	}

	var clientRequest domain.ClientResponse
	if err := json.NewDecoder(resp.Body).Decode(&clientRequest); err != nil {
		return nil, errors.New("failed parsing request body")
	}
	return &clientRequest, nil
}
