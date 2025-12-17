package mappers

import (
	"invest-mate/internal/models/domain"
	"invest-mate/internal/models/dto"
	"invest-mate/internal/models/entity"
)

func ShareFromDtoMapper(dto dto.Share) domain.Share {
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
	}
}

func ShareToRepositoryMapper(dto []domain.Share) []entity.Share {
	res := make([]entity.Share, len(dto))
	for i, b := range dto {
		res[i] = entity.Share{
			Figi:                  b.Figi,
			Ticker:                b.Ticker,
			PositionUid:           b.PositionUid,
			ClassCode:             b.ClassCode,
			Isin:                  b.Isin,
			Lot:                   b.Lot,
			Currency:              b.Currency,
			Klong:                 b.Klong,
			Kshort:                b.Kshort,
			Dlong:                 b.Dlong,
			Dshort:                b.Dshort,
			DlongMin:              b.DlongMin,
			DshortMin:             b.DshortMin,
			Exchange:              b.Exchange,
			RealExchange:          b.RealExchange,
			ShortEnabledFlag:      b.ShortEnabledFlag,
			Name:                  b.Name,
			CountryOfRiskName:     b.CountryOfRiskName,
			TradingStatus:         b.TradingStatus,
			OtcFlag:               b.OtcFlag,
			BuyAvailableFlag:      b.BuyAvailableFlag,
			SellAvailableFlag:     b.SellAvailableFlag,
			MinPriceIncrement:     b.MinPriceIncrement,
			ApiTradeAvailableFlag: b.ApiTradeAvailableFlag,
			Uid:                   b.Uid,
			AssetUid:              b.AssetUid,
			ForIisFlag:            b.ForIisFlag,
			ForQualInvestorFlag:   b.ForQualInvestorFlag,
			WeekendFlag:           b.WeekendFlag,
			BlockedTcaFlag:        b.BlockedTcaFlag,
			LiquidityFlag:         b.LiquidityFlag,
			First1minCandleDate:   b.First1minCandleDate,
			First1dayCandleDate:   b.First1dayCandleDate,
			DivYieldFlag:          b.DivYieldFlag,
			IpoDate:               b.IpoDate,
			IssueSize:             b.IssueSize,
			IssueSizePlan:         b.IssueSizePlan,
			Nominal:               b.Nominal,
			Sector:                b.Sector,
			ShareType:             b.ShareType,
		}
	}
	return res
}
