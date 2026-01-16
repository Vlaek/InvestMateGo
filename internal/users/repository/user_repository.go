package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"invest-mate/internal/users/mappers"
	"invest-mate/internal/users/models"
	"invest-mate/internal/users/models/domain"
	"invest-mate/internal/users/models/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByField(ctx context.Context, fieldName string, fieldValue string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) (bool, error)
	GetList(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// Создание нового репозитория
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Создание нового пользователя в БД
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	var count int64
	r.db.WithContext(ctx).Model(&entity.User{}).
		Where("email = ?", user.Email).
		Count(&count)
	if count > 0 {
		return models.ErrEmailAlreadyExists
	}

	entityUser := mappers.FromDomainToEntity(user)

	if err := r.db.WithContext(ctx).Create(&entityUser).Error; err != nil {
		return err
	}

	user.ID = entityUser.ID
	user.CreatedAt = entityUser.CreatedAt
	user.UpdatedAt = entityUser.UpdatedAt

	return nil
}

// Найти пользователя по полю в БД
func (r *userRepository) FindByField(ctx context.Context, fieldName string, fieldValue string) (*domain.User, error) {
	var entityUser entity.User
	var query string = fmt.Sprintf("%s = ?", fieldName)
	err := r.db.WithContext(ctx).First(&entityUser, query, fieldValue).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return mappers.FromEntityToDomain(entityUser), nil
}

// Обновить пользователя в БД
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	entityUser := mappers.FromDomainToEntity(user)

	if user.Email != entityUser.Email {
		var count int64
		r.db.WithContext(ctx).Model(&entity.User{}).
			Where("email = ? AND id != ?", user.Email, user.ID).
			Count(&count)
		if count > 0 {
			return models.ErrEmailAlreadyExists
		}
	}

	return r.db.WithContext(ctx).Save(&entityUser).Error
}

// Удаление пользователя из БД
func (r *userRepository) Delete(ctx context.Context, id string) (bool, error) {
	result := r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id)

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

// Получить список пользователей в БД
func (r *userRepository) GetList(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var entityUsers []entity.User

	query := r.db.WithContext(ctx).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&entityUsers).Error
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(entityUsers))
	for i, eu := range entityUsers {
		users[i] = mappers.FromEntityToDomain(eu)
	}

	return users, nil
}
