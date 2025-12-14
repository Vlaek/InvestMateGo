package storage

import (
	models "invest-mate/internal/models"
	repositoryModels "invest-mate/internal/repository/models"
)

func BondToRepositoryMapper(dto []models.Bond) []repositoryModels.Bond {
	res := make([]repositoryModels.Bond, len(dto))
	for i, b := range dto {
		res[i] = repositoryModels.Bond{
			Uid:                   b.Uid,
			Figi:                  b.Figi,
			Ticker:                b.Ticker,
			Name:                  b.Name,
			Currency:              b.Currency,
			PositionUid:           b.PositionUid,
			ClassCode:             b.ClassCode,
			Isin:                  b.Isin,
			Lot:                   b.Lot,
			Klong:                 b.Klong,
			Kshort:                b.Kshort,
			Dlong:                 b.Dlong,
			Dshort:                b.Dshort,
			DlongMin:              b.DlongMin,
			DshortMin:             b.DshortMin,
			Exchange:              b.Exchange,
			RealExchange:          b.RealExchange,
			ShortEnabledFlag:      b.ShortEnabledFlag,
			CountryOfRiskName:     b.CountryOfRiskName,
			TradingStatus:         b.TradingStatus,
			OtcFlag:               b.OtcFlag,
			BuyAvailableFlag:      b.BuyAvailableFlag,
			SellAvailableFlag:     b.SellAvailableFlag,
			MinPriceIncrement:     b.MinPriceIncrement,
			ApiTradeAvailableFlag: b.ApiTradeAvailableFlag,
			AssetUid:              b.AssetUid,
			ForIisFlag:            b.ForIisFlag,
			ForQualInvestorFlag:   b.ForQualInvestorFlag,
			WeekendFlag:           b.WeekendFlag,
			BlockedTcaFlag:        b.BlockedTcaFlag,
			LiquidityFlag:         b.LiquidityFlag,
			First1minCandleDate:   b.First1minCandleDate,
			First1dayCandleDate:   b.First1dayCandleDate,
			AciValue:              b.AciValue,
			AmortizationFlag:      b.AmortizationFlag,
			CountryOfRisk:         b.CountryOfRisk,
			CouponQuantityPerYear: b.CouponQuantityPerYear,
			FloatingCouponFlag:    b.FloatingCouponFlag,
			InitialNominal:        b.InitialNominal,
			IssueKind:             b.IssueKind,
			IssueSize:             b.IssueSize,
			IssueSizePlan:         b.IssueSizePlan,
			MaturityDate:          b.MaturityDate,
			Nominal:               b.Nominal,
			PerpetualFlag:         b.PerpetualFlag,
			PlacementDate:         b.PlacementDate,
			PlacementPrice:        b.PlacementPrice,
			RiskLevel:             b.RiskLevel,
			Sector:                b.Sector,
			StateRegDate:          b.StateRegDate,
			SubordinatedFlag:      b.SubordinatedFlag,
			BondType:              b.BondType,
		}
	}
	return res
}
