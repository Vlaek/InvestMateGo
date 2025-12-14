package storage

import (
	models "invest-mate/internal/models"
	repositoryModels "invest-mate/internal/repository/models"
)

func ShareToRepositoryMapper(dto []models.Share) []repositoryModels.Share {
	res := make([]repositoryModels.Share, len(dto))
	for i, b := range dto {
		res[i] = repositoryModels.Share{
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
