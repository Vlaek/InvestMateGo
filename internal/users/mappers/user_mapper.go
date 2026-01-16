package mappers

import (
	"invest-mate/internal/users/models/domain"
	"invest-mate/internal/users/models/entity"
)

func FromEntityToDomain(entity entity.User) *domain.User {
	return &domain.User{
		ID:           entity.ID,
		Email:        entity.Email,
		Username:     entity.Username,
		PasswordHash: entity.PasswordHash,
		Role:         entity.Role,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

func FromEntityToDomainSlice(entitySlice []entity.User) []*domain.User {
	domainSlice := make([]*domain.User, len(entitySlice))

	for index, dto := range entitySlice {
		domainSlice[index] = FromEntityToDomain(dto)
	}

	return domainSlice
}

func FromDomainToEntity(domain *domain.User) entity.User {
	return entity.User{
		ID:           domain.ID,
		Email:        domain.Email,
		Username:     domain.Username,
		PasswordHash: domain.PasswordHash,
		Role:         domain.Role,
		CreatedAt:    domain.CreatedAt,
		UpdatedAt:    domain.UpdatedAt,
	}
}

func FromDomainToEntitySlice(domainSlice []*domain.User) []entity.User {
	entitySlice := make([]entity.User, len(domainSlice))

	for index, domain := range domainSlice {
		entitySlice[index] = FromDomainToEntity(domain)
	}

	return entitySlice
}
