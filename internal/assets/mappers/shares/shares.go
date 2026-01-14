package shares

import (
	"invest-mate/internal/assets/models/domain"
	"invest-mate/internal/assets/models/dto"
	"invest-mate/internal/assets/models/entity"
	"invest-mate/internal/shared/models"
)

func FromDtoToDomain(dto dto.Share) domain.Share {
	return domain.Share{
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
		DivYieldFlag:          dto.DivYieldFlag,
		IpoDate:               dto.IpoDate,
		IssueSize:             dto.IssueSize,
		IssueSizePlan:         dto.IssueSizePlan,
		Nominal:               dto.Nominal.ToFloat(),
		Sector:                dto.Sector,
		ShareType:             dto.ShareType,

		InstrumentType: models.InstrumentTypeBond,
	}
}

func FromDtoToDomainSlice(dtoSlice []dto.Share) []domain.Share {
	domainSlice := make([]domain.Share, len(dtoSlice))

	for index, dto := range dtoSlice {
		domainSlice[index] = FromDtoToDomain(dto)
	}

	return domainSlice
}

func FromDomainToEntity(domain domain.Share) entity.Share {
	return entity.Share{
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
		DivYieldFlag:          domain.DivYieldFlag,
		IpoDate:               domain.IpoDate,
		IssueSize:             domain.IssueSize,
		IssueSizePlan:         domain.IssueSizePlan,
		Nominal:               domain.Nominal,
		Sector:                domain.Sector,
		ShareType:             domain.ShareType,

		InstrumentType: models.InstrumentTypeBond,
	}
}

func FromDomainToEntitySlice(domainSlice []domain.Share) []entity.Share {
	entitySlice := make([]entity.Share, len(domainSlice))

	for index, domain := range domainSlice {
		entitySlice[index] = FromDomainToEntity(domain)
	}

	return entitySlice
}

func FromEntityToDomain(entity entity.Share) domain.Share {
	return domain.Share{
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
		DivYieldFlag:          entity.DivYieldFlag,
		IpoDate:               entity.IpoDate,
		IssueSize:             entity.IssueSize,
		IssueSizePlan:         entity.IssueSizePlan,
		Nominal:               entity.Nominal,
		Sector:                entity.Sector,
		ShareType:             entity.ShareType,

		InstrumentType: models.InstrumentTypeBond,
	}
}

func FromEntityToDomainSlice(entitySlice []entity.Share) []domain.Share {
	domainSlice := make([]domain.Share, len(entitySlice))

	for index, entity := range entitySlice {
		domainSlice[index] = FromEntityToDomain(entity)
	}

	return domainSlice
}
