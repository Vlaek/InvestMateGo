package assets

import (
	"invest-mate/internal/assets/models/domain"
	"invest-mate/internal/assets/models/dto"
	"invest-mate/internal/assets/models/entity"
	"invest-mate/internal/shared/models"
)

func FromDtoToDomain(marker dto.Marker) domain.Asset {
	switch v := marker.(type) {
	case dto.Bond:
		{
			return domain.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeBond,
			}
		}
	case dto.Share:
		{
			return domain.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeShare,
			}
		}
	case dto.Etf:
		{
			return domain.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeETF,
			}
		}
	case dto.Currency:
		{
			return domain.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeCurrency,
			}
		}
	default:
		{
			return domain.Asset{}
		}
	}
}

func FromDtoToDomainSlice[T dto.Marker](dtoSlice []T) []domain.Asset {
	domainSlice := make([]domain.Asset, len(dtoSlice))

	for index, dto := range dtoSlice {
		domainSlice[index] = FromDtoToDomain(dto)
	}

	return domainSlice
}

func FromDomainToEntity(marker domain.Marker) entity.Asset {
	switch v := marker.(type) {
	case domain.Bond:
		{
			return entity.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeBond,
			}
		}
	case domain.Share:
		{
			return entity.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeShare,
			}
		}
	case domain.Etf:
		{
			return entity.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeETF,
			}
		}
	case domain.Currency:
		{
			return entity.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeCurrency,
			}
		}
	default:
		{
			return entity.Asset{}
		}
	}
}

func FromDomainToEntitySlice[T domain.Marker](domainSlice []T) []entity.Asset {
	entitySlice := make([]entity.Asset, len(domainSlice))

	for index, domain := range domainSlice {
		entitySlice[index] = FromDomainToEntity(domain)
	}

	return entitySlice
}

func FromEntityToDomain(marker entity.Marker) domain.Asset {
	switch v := marker.(type) {
	case entity.Bond:
		{
			return domain.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeBond,
			}
		}
	case entity.Share:
		{
			return domain.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeShare,
			}
		}
	case entity.Etf:
		{
			return domain.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeETF,
			}
		}
	case entity.Currency:
		{
			return domain.Asset{
				Uid: v.Uid,

				InstrumentType: models.InstrumentTypeCurrency,
			}
		}
	default:
		{
			return domain.Asset{}
		}
	}
}

func FromEntityToDomainSlice[T entity.Marker](entitySlice []T) []domain.Asset {
	domainSlice := make([]domain.Asset, len(entitySlice))

	for index, entity := range entitySlice {
		domainSlice[index] = FromEntityToDomain(entity)
	}

	return domainSlice
}
