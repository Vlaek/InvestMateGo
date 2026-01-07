package etfs

import (
	"invest-mate/internal/models/domain"
	"invest-mate/internal/models/dto"
	"invest-mate/internal/models/entity"
)

func FromDtoToDomain(dto dto.Etf) domain.Etf {
	return domain.Etf{
		Figi:                  dto.Figi,
		Ticker:                dto.Ticker,
		PositionUid:           dto.PositionUid,
		ClassCode:             dto.ClassCode,
		Isin:                  dto.Isin,
		Lot:                   dto.Lot,
		Currency:              dto.Currency,
		Klong:                 dto.Klong.ToFloat(),
		Kshort:                dto.Kshort.ToFloat(),
		Dlong:                 dto.Dlong.ToFloat(),
		Dshort:                dto.Dshort.ToFloat(),
		DlongMin:              dto.DlongMin.ToFloat(),
		DshortMin:             dto.DshortMin.ToFloat(),
		Exchange:              dto.Exchange,
		RealExchange:          dto.RealExchange,
		ShortEnabledFlag:      dto.ShortEnabledFlag,
		Name:                  dto.Name,
		CountryOfRiskName:     dto.CountryOfRiskName,
		TradingStatus:         dto.TradingStatus,
		OtcFlag:               dto.OtcFlag,
		BuyAvailableFlag:      dto.BuyAvailableFlag,
		SellAvailableFlag:     dto.SellAvailableFlag,
		MinPriceIncrement:     dto.MinPriceIncrement.ToFloat(),
		ApiTradeAvailableFlag: dto.ApiTradeAvailableFlag,
		Uid:                   dto.Uid,
		AssetUid:              dto.AssetUid,
		ForIisFlag:            dto.ForIisFlag,
		ForQualInvestorFlag:   dto.ForQualInvestorFlag,
		WeekendFlag:           dto.WeekendFlag,
		BlockedTcaFlag:        dto.BlockedTcaFlag,
		LiquidityFlag:         dto.LiquidityFlag,
		First1minCandleDate:   dto.First1minCandleDate,
		First1dayCandleDate:   dto.First1dayCandleDate,
		FixedCommission:       dto.FixedCommission.ToFloat(),
		FocusType:             dto.FocusType,
		ReleasedDate:          dto.ReleasedDate,
		NumShares:             dto.NumShares.ToFloat(),
		Sector:                dto.Sector,
		RebalancingFreq:       dto.RebalancingFreq,
	}
}

func FromDtoToDomainSlice(dtoSlice []dto.Etf) []domain.Etf {
	domainSlice := make([]domain.Etf, len(dtoSlice))

	for index, dto := range dtoSlice {
		domainSlice[index] = FromDtoToDomain(dto)
	}

	return domainSlice
}

func FromDomainToEntity(domain domain.Etf) entity.Etf {
	return entity.Etf{
		Figi:                  domain.Figi,
		Ticker:                domain.Ticker,
		PositionUid:           domain.PositionUid,
		ClassCode:             domain.ClassCode,
		Isin:                  domain.Isin,
		Lot:                   domain.Lot,
		Currency:              domain.Currency,
		Klong:                 domain.Klong,
		Kshort:                domain.Kshort,
		Dlong:                 domain.Dlong,
		Dshort:                domain.Dshort,
		DlongMin:              domain.DlongMin,
		DshortMin:             domain.DshortMin,
		Exchange:              domain.Exchange,
		RealExchange:          domain.RealExchange,
		ShortEnabledFlag:      domain.ShortEnabledFlag,
		Name:                  domain.Name,
		CountryOfRiskName:     domain.CountryOfRiskName,
		TradingStatus:         domain.TradingStatus,
		OtcFlag:               domain.OtcFlag,
		BuyAvailableFlag:      domain.BuyAvailableFlag,
		SellAvailableFlag:     domain.SellAvailableFlag,
		MinPriceIncrement:     domain.MinPriceIncrement,
		ApiTradeAvailableFlag: domain.ApiTradeAvailableFlag,
		Uid:                   domain.Uid,
		AssetUid:              domain.AssetUid,
		ForIisFlag:            domain.ForIisFlag,
		ForQualInvestorFlag:   domain.ForQualInvestorFlag,
		WeekendFlag:           domain.WeekendFlag,
		BlockedTcaFlag:        domain.BlockedTcaFlag,
		LiquidityFlag:         domain.LiquidityFlag,
		First1minCandleDate:   domain.First1minCandleDate,
		First1dayCandleDate:   domain.First1dayCandleDate,
		FixedCommission:       domain.FixedCommission,
		FocusType:             domain.FocusType,
		ReleasedDate:          domain.ReleasedDate,
		NumShares:             domain.NumShares,
		Sector:                domain.Sector,
		RebalancingFreq:       domain.RebalancingFreq,
	}
}

func FromDomainToEntitySlice(domainSlice []domain.Etf) []entity.Etf {
	entitySlice := make([]entity.Etf, len(domainSlice))

	for index, domain := range domainSlice {
		entitySlice[index] = FromDomainToEntity(domain)
	}

	return entitySlice
}

func FromEntityToDomain(entity entity.Etf) domain.Etf {
	return domain.Etf{
		Figi:                  entity.Figi,
		Ticker:                entity.Ticker,
		PositionUid:           entity.PositionUid,
		ClassCode:             entity.ClassCode,
		Isin:                  entity.Isin,
		Lot:                   entity.Lot,
		Currency:              entity.Currency,
		Klong:                 entity.Klong,
		Kshort:                entity.Kshort,
		Dlong:                 entity.Dlong,
		Dshort:                entity.Dshort,
		DlongMin:              entity.DlongMin,
		DshortMin:             entity.DshortMin,
		Exchange:              entity.Exchange,
		RealExchange:          entity.RealExchange,
		ShortEnabledFlag:      entity.ShortEnabledFlag,
		Name:                  entity.Name,
		CountryOfRiskName:     entity.CountryOfRiskName,
		TradingStatus:         entity.TradingStatus,
		OtcFlag:               entity.OtcFlag,
		BuyAvailableFlag:      entity.BuyAvailableFlag,
		SellAvailableFlag:     entity.SellAvailableFlag,
		MinPriceIncrement:     entity.MinPriceIncrement,
		ApiTradeAvailableFlag: entity.ApiTradeAvailableFlag,
		Uid:                   entity.Uid,
		AssetUid:              entity.AssetUid,
		ForIisFlag:            entity.ForIisFlag,
		ForQualInvestorFlag:   entity.ForQualInvestorFlag,
		WeekendFlag:           entity.WeekendFlag,
		BlockedTcaFlag:        entity.BlockedTcaFlag,
		LiquidityFlag:         entity.LiquidityFlag,
		First1minCandleDate:   entity.First1minCandleDate,
		First1dayCandleDate:   entity.First1dayCandleDate,
		FixedCommission:       entity.FixedCommission,
		FocusType:             entity.FocusType,
		ReleasedDate:          entity.ReleasedDate,
		NumShares:             entity.NumShares,
		Sector:                entity.Sector,
		RebalancingFreq:       entity.RebalancingFreq,
	}
}

func FromEntityToDomainSlice(entitySlice []entity.Etf) []domain.Etf {
	domainSlice := make([]domain.Etf, len(entitySlice))

	for index, entity := range entitySlice {
		domainSlice[index] = FromEntityToDomain(entity)
	}

	return domainSlice
}
