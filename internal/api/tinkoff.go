package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"invest-mate/internal/config"
	"invest-mate/internal/mappers/bonds"
	"invest-mate/internal/mappers/currencies"
	"invest-mate/internal/mappers/etfs"
	"invest-mate/internal/mappers/shares"
	"invest-mate/internal/models/domain"
	"invest-mate/internal/models/dto"
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

func (c *TinkoffClient) doRequest(ctx context.Context, method, endpoint string, body any) (*http.Response, error) {
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

func handleAPIError(resp *http.Response, endpoint string) error {
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

func fetchInstruments[T dto.Marker, M domain.Marker](
	ctx context.Context,
	client *TinkoffClient,
	endpoint string,
	mapper func(dto T) M,
) ([]M, error) {
	body := map[string]string{
		"instrumentStatus": "INSTRUMENT_STATUS_BASE",
	}

	resp, err := client.doRequest(ctx, "POST", endpoint, body)

	if err != nil {
		return nil, fmt.Errorf("request %s: %w", endpoint, err)
	}

	defer resp.Body.Close()

	logger.InfoLog("%s API response status: %d", endpoint, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, handleAPIError(resp, endpoint)
	}

	var dtoResponse struct {
		Instruments []T `json:"instruments"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(resp.Body).Decode(&dtoResponse); err != nil {
		logger.ErrorLog("Failed to decode JSON for %s. Body start: %s",
			endpoint, string(bodyBytes[:min(500, len(bodyBytes))]))

		return nil, fmt.Errorf("decode DTO response for %s: %w", endpoint, err)
	}

	logger.InfoLog("Successfully parsed %d DTO instruments from %s",
		len(dtoResponse.Instruments), endpoint)

	instruments := make([]M, len(dtoResponse.Instruments))

	for i, dtoItem := range dtoResponse.Instruments {
		instruments[i] = mapper(dtoItem)
	}

	logger.InfoLog("Mapping complete for %s: total: %d", endpoint, len(instruments))

	return instruments, nil
}

func GetBonds(ctx context.Context) ([]domain.Bond, error) {
	client := NewTinkoffClient()

	return fetchInstruments(
		ctx,
		client,
		"tinkoff.public.invest.api.contract.v1.InstrumentsService/Bonds",
		bonds.FromDtoToDomain,
	)
}

func GetShares(ctx context.Context) ([]domain.Share, error) {
	client := NewTinkoffClient()

	return fetchInstruments(
		ctx,
		client,
		"tinkoff.public.invest.api.contract.v1.InstrumentsService/Shares",
		shares.FromDtoToDomain,
	)
}

func GetEtfs(ctx context.Context) ([]domain.Etf, error) {
	client := NewTinkoffClient()

	return fetchInstruments(
		ctx,
		client,
		"tinkoff.public.invest.api.contract.v1.InstrumentsService/Etfs",
		etfs.FromDtoToDomain,
	)
}

func GetCurrencies(ctx context.Context) ([]domain.Currency, error) {
	client := NewTinkoffClient()

	return fetchInstruments(
		ctx,
		client,
		"tinkoff.public.invest.api.contract.v1.InstrumentsService/Currencies",
		currencies.FromDtoToDomain,
	)
}
