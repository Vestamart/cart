package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vestamart/cart/internal/domain"
	"github.com/vestamart/cart/internal/localErr"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	url        string
	token      string
	limiter    <-chan time.Time
	ticker     *time.Ticker
}

func NewClient(url, token string) *Client {
	ticker := time.NewTicker(time.Millisecond * 100)
	return &Client{
		httpClient: &http.Client{Timeout: time.Second * 5},
		url:        url,
		token:      token,
		limiter:    ticker.C,
		ticker:     ticker,
	}
}

type request struct {
	Token string `json:"token"`
	SKU   int64  `json:"sku"`
}

func (c *Client) Close() {
	c.ticker.Stop()
	c.httpClient.CloseIdleConnections()
}

func (c *Client) ExistItem(ctx context.Context, sku int64) error {
	select {
	case <-c.limiter:
	case <-ctx.Done():
		return ctx.Err()
	}

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
		return localErr.ErrSkuNotExist
	} else if resp.StatusCode != http.StatusOK {
		return errors.New("error exist item")
	}

	return nil
}

func (c *Client) GetProduct(ctx context.Context, sku int64) (*domain.ProductServiceResponse, error) {
	select {
	case <-c.limiter:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

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
