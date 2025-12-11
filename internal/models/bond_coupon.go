package models

type BondCoupon struct {
    Figi         string  `json:"figi"`
    CouponDate   string  `json:"coupon_date"`
    CouponNumber int     `json:"coupon_number"`
}