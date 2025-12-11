package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"invest-mate/config"
	"invest-mate/internal/api/dto"
	"invest-mate/internal/api/mappers"
	"invest-mate/internal/models"
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

func (c *TinkoffClient) doRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
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

	// ДЕТАЛЬНОЕ логирование для отладки
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

func GetBonds(ctx context.Context) ([]models.Bond, error) {
	client := NewTinkoffClient()

	body := map[string]string{
		"instrumentStatus": "INSTRUMENT_STATUS_BASE",
	}

	resp, err := client.doRequest(ctx, "POST",
		"tinkoff.public.invest.api.contract.v1.InstrumentsService/Bonds",
		body)
	if err != nil {
		return nil, fmt.Errorf("request bonds: %w", err)
	}
	defer resp.Body.Close()

	logger.InfoLog("Bonds API response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
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

		return nil, fmt.Errorf("bonds API error: %s", errorMsg)
	}

	var dtoResponse struct {
		Instruments []dto.BondDTO `json:"instruments"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(resp.Body).Decode(&dtoResponse); err != nil {
		logger.ErrorLog("Failed to decode JSON. Body start: %s",
			string(bodyBytes[:min(500, len(bodyBytes))]))
		return nil, fmt.Errorf("decode DTO response: %w", err)
	}

	logger.InfoLog("Successfully parsed %d DTO bonds", len(dtoResponse.Instruments))

	bonds := make([]models.Bond, 0, len(dtoResponse.Instruments))

	for _, dtoItem := range dtoResponse.Instruments {
		bond := mappers.BondFromDTO(dtoItem)
		bonds = append(bonds, bond)
	}

	logger.InfoLog("Mapping complete: total: %d", len(bonds))

	return bonds, nil
}
