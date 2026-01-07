package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"invest-mate/internal/common/config"
	"invest-mate/pkg/logger"
)

type TinkoffClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewTinkoffClient() *TinkoffClient {
	cfg := config.Get()

	return &TinkoffClient{
		baseURL: "https://invest-public-api.tbank.ru/rest/",
		token:   cfg.TinkoffToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *TinkoffClient) DoRequest(ctx context.Context, method, endpoint string, body any) (*http.Response, error) {
	url := c.baseURL + endpoint
	var reqBody []byte

	if body != nil {
		var err error

		reqBody, err = json.Marshal(body)

		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	logger.InfoLog("Making %s request to: %s", method, url)

	if body != nil {
		logger.InfoLog("Request body: %s", string(reqBody))
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("do request to %s: %w", url, err)
	}

	return resp, nil
}

func (c *TinkoffClient) HandleAPIError(resp *http.Response, endpoint string) error {
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)
	errorMsg := fmt.Sprintf("status %d", resp.StatusCode)

	if len(bodyStr) > 0 {
		displayLen := 200

		if len(bodyStr) < displayLen {
			displayLen = len(bodyStr)
		}

		errorMsg = fmt.Sprintf("status %d: %s", resp.StatusCode, bodyStr[:displayLen])
	}

	return fmt.Errorf("%s API error: %s", endpoint, errorMsg)
}
