package domain

import "invest-mate/internal/shared/models"

type Asset struct {
	Uid            string                `json:"uid"`
	InstrumentType models.InstrumentType `json:"instrumentType"`
}
