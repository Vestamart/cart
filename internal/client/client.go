package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vestamart/homework/internal/domain"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	url        string
	token      string
}

func NewClient(url, token string) *Client {
	return &Client{
		httpClient: &http.Client{},
		url:        url,
		token:      token,
	}
}

type request struct {
	Token string `json:"token"`
	SKU   int64  `json:"sku"`
}

func (c *Client) ExistItem(ctx context.Context, sku int64) error {
	jsonBody, err := json.Marshal(request{Token: c.token, SKU: sku})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.ErrSkuNotExist
	} else if resp.StatusCode != http.StatusOK {
		return errors.New("error exist item")
	}

	return nil
}

func (c *Client) GetProduct(ctx context.Context, sku int64) (*domain.ProductServiceResponse, error) {
	jsonBody, err := json.Marshal(request{Token: c.token, SKU: sku})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error get product")
	}

	var clientResponse domain.ProductServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&clientResponse); err != nil {
		return nil, errors.New("failed parsing request body")
	}
	return &clientResponse, nil
}
