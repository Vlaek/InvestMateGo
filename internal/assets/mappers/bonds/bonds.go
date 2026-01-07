package bonds

import (
	"invest-mate/internal/assets/models/domain"
	"invest-mate/internal/assets/models/dto"
	"invest-mate/internal/assets/models/entity"
)

func FromDtoToDomain(dto dto.Bond) domain.Bond {
	return domain.Bond{
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
		AciValue:              dto.AciValue.ToFloat(),
		AmortizationFlag:      dto.AmortizationFlag,
		CountryOfRisk:         dto.CountryOfRisk,
		CouponQuantityPerYear: dto.CouponQuantityPerYear,
		FloatingCouponFlag:    dto.FloatingCouponFlag,
		InitialNominal:        dto.InitialNominal.ToFloat(),
		IssueKind:             dto.IssueKind,
		IssueSize:             dto.IssueSize,
		IssueSizePlan:         dto.IssueSizePlan,
		MaturityDate:          dto.MaturityDate,
		Nominal:               dto.Nominal.ToFloat(),
		PerpetualFlag:         dto.PerpetualFlag,
		PlacementDate:         dto.PlacementDate,
		PlacementPrice:        dto.PlacementPrice.ToFloat(),
		RiskLevel:             dto.RiskLevel.ToString(),
		Sector:                dto.Sector,
		StateRegDate:          dto.StateRegDate,
		SubordinatedFlag:      dto.SubordinatedFlag,
		BondType:              dto.BondType,
	}
}

func FromDtoToDomainSlice(dtoSlice []dto.Bond) []domain.Bond {
	domainSlice := make([]domain.Bond, len(dtoSlice))

	for index, dto := range dtoSlice {
		domainSlice[index] = FromDtoToDomain(dto)
	}

	return domainSlice
}

func FromDomainToEntity(domain domain.Bond) entity.Bond {
	return entity.Bond{
		Uid:                   domain.Uid,
		Figi:                  domain.Figi,
		Ticker:                domain.Ticker,
		Name:                  domain.Name,
		Currency:              domain.Currency,
		PositionUid:           domain.PositionUid,
		ClassCode:             domain.ClassCode,
		Isin:                  domain.Isin,
		Lot:                   domain.Lot,
		Klong:                 domain.Klong,
		Kshort:                domain.Kshort,
		Dlong:                 domain.Dlong,
		Dshort:                domain.Dshort,
		DlongMin:              domain.DlongMin,
		DshortMin:             domain.DshortMin,
		Exchange:              domain.Exchange,
		RealExchange:          domain.RealExchange,
		ShortEnabledFlag:      domain.ShortEnabledFlag,
		CountryOfRiskName:     domain.CountryOfRiskName,
		TradingStatus:         domain.TradingStatus,
		OtcFlag:               domain.OtcFlag,
		BuyAvailableFlag:      domain.BuyAvailableFlag,
		SellAvailableFlag:     domain.SellAvailableFlag,
		MinPriceIncrement:     domain.MinPriceIncrement,
		ApiTradeAvailableFlag: domain.ApiTradeAvailableFlag,
		AssetUid:              domain.AssetUid,
		ForIisFlag:            domain.ForIisFlag,
		ForQualInvestorFlag:   domain.ForQualInvestorFlag,
		WeekendFlag:           domain.WeekendFlag,
		BlockedTcaFlag:        domain.BlockedTcaFlag,
		LiquidityFlag:         domain.LiquidityFlag,
		First1minCandleDate:   domain.First1minCandleDate,
		First1dayCandleDate:   domain.First1dayCandleDate,
		AciValue:              domain.AciValue,
		AmortizationFlag:      domain.AmortizationFlag,
		CountryOfRisk:         domain.CountryOfRisk,
		CouponQuantityPerYear: domain.CouponQuantityPerYear,
		FloatingCouponFlag:    domain.FloatingCouponFlag,
		InitialNominal:        domain.InitialNominal,
		IssueKind:             domain.IssueKind,
		IssueSize:             domain.IssueSize,
		IssueSizePlan:         domain.IssueSizePlan,
		MaturityDate:          domain.MaturityDate,
		Nominal:               domain.Nominal,
		PerpetualFlag:         domain.PerpetualFlag,
		PlacementDate:         domain.PlacementDate,
		PlacementPrice:        domain.PlacementPrice,
		RiskLevel:             domain.RiskLevel,
		Sector:                domain.Sector,
		StateRegDate:          domain.StateRegDate,
		SubordinatedFlag:      domain.SubordinatedFlag,
		BondType:              domain.BondType,
	}
}

func FromDomainToEntitySlice(domainSlice []domain.Bond) []entity.Bond {
	entitySlice := make([]entity.Bond, len(domainSlice))

	for index, domain := range domainSlice {
		entitySlice[index] = FromDomainToEntity(domain)
	}

	return entitySlice
}

func FromEntityToDomain(entity entity.Bond) domain.Bond {
	return domain.Bond{
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
		AciValue:              entity.AciValue,
		AmortizationFlag:      entity.AmortizationFlag,
		CountryOfRisk:         entity.CountryOfRisk,
		CouponQuantityPerYear: entity.CouponQuantityPerYear,
		FloatingCouponFlag:    entity.FloatingCouponFlag,
		InitialNominal:        entity.InitialNominal,
		IssueKind:             entity.IssueKind,
		IssueSize:             entity.IssueSize,
		IssueSizePlan:         entity.IssueSizePlan,
		MaturityDate:          entity.MaturityDate,
		Nominal:               entity.Nominal,
		PerpetualFlag:         entity.PerpetualFlag,
		PlacementDate:         entity.PlacementDate,
		PlacementPrice:        entity.PlacementPrice,
		RiskLevel:             entity.RiskLevel,
		Sector:                entity.Sector,
		StateRegDate:          entity.StateRegDate,
		SubordinatedFlag:      entity.SubordinatedFlag,
		BondType:              entity.BondType,
	}
}

func FromEntityToDomainSlice(entitySlice []entity.Bond) []domain.Bond {
	domainSlice := make([]domain.Bond, len(entitySlice))

	for index, entity := range entitySlice {
		domainSlice[index] = FromEntityToDomain(entity)
	}

	return domainSlice
}
