package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"invest-mate/internal/assets/mappers/bonds"
	"invest-mate/internal/assets/mappers/currencies"
	"invest-mate/internal/assets/mappers/etfs"
	"invest-mate/internal/assets/mappers/shares"
	"invest-mate/internal/assets/models/domain"
	"invest-mate/internal/assets/models/dto"
	commonAPI "invest-mate/internal/common/api"
	"invest-mate/pkg/logger"
)

func fetchInstruments[T dto.Marker, M domain.Marker](
	ctx context.Context,
	client *commonAPI.TinkoffClient,
	endpoint string,
	mapper func(dto T) M,
) ([]M, error) {
	body := map[string]string{
		"instrumentStatus": "INSTRUMENT_STATUS_BASE",
	}

	resp, err := client.DoRequest(ctx, "POST", endpoint, body)

	if err != nil {
		return nil, fmt.Errorf("request %s: %w", endpoint, err)
	}

	defer resp.Body.Close()

	logger.InfoLog("%s API response status: %d", endpoint, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, client.HandleAPIError(resp, endpoint)
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
	client := commonAPI.NewTinkoffClient()

	return fetchInstruments(
		ctx,
		client,
		"tinkoff.public.invest.api.contract.v1.InstrumentsService/Bonds",
		bonds.FromDtoToDomain,
	)
}

func GetShares(ctx context.Context) ([]domain.Share, error) {
	client := commonAPI.NewTinkoffClient()

	return fetchInstruments(
		ctx,
		client,
		"tinkoff.public.invest.api.contract.v1.InstrumentsService/Shares",
		shares.FromDtoToDomain,
	)
}

func GetEtfs(ctx context.Context) ([]domain.Etf, error) {
	client := commonAPI.NewTinkoffClient()

	return fetchInstruments(
		ctx,
		client,
		"tinkoff.public.invest.api.contract.v1.InstrumentsService/Etfs",
		etfs.FromDtoToDomain,
	)
}

func GetCurrencies(ctx context.Context) ([]domain.Currency, error) {
	client := commonAPI.NewTinkoffClient()

	return fetchInstruments(
		ctx,
		client,
		"tinkoff.public.invest.api.contract.v1.InstrumentsService/Currencies",
		currencies.FromDtoToDomain,
	)
}
