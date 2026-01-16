package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"invest-mate/internal/users/models/domain"
	"invest-mate/internal/users/models/entity"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUsernameTaken      = errors.New("username already taken")
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	// Проверяем уникальность email
	var count int64
	r.db.WithContext(ctx).Model(&entity.User{}).
		Where("email = ?", user.Email).
		Count(&count)
	if count > 0 {
		return ErrEmailAlreadyExists
	}

	// Проверяем уникальность username
	r.db.WithContext(ctx).Model(&entity.User{}).
		Where("username = ?", user.Username).
		Count(&count)
	if count > 0 {
		return ErrUsernameTaken
	}

	// Преобразуем доменную модель в сущность
	entityUser := r.toEntity(user)

	// Создаем пользователя
	if err := r.db.WithContext(ctx).Create(&entityUser).Error; err != nil {
		return err
	}

	// Обновляем ID в доменной модели
	user.ID = entityUser.ID
	user.CreatedAt = entityUser.CreatedAt
	user.UpdatedAt = entityUser.UpdatedAt

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var entityUser entity.User
	err := r.db.WithContext(ctx).First(&entityUser, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return r.toDomain(&entityUser), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var entityUser entity.User
	err := r.db.WithContext(ctx).First(&entityUser, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return r.toDomain(&entityUser), nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	entityUser := r.toEntity(user)

	// Проверяем уникальность email (если изменился)
	if user.Email != entityUser.Email {
		var count int64
		r.db.WithContext(ctx).Model(&entity.User{}).
			Where("email = ? AND id != ?", user.Email, user.ID).
			Count(&count)
		if count > 0 {
			return ErrEmailAlreadyExists
		}
	}

	// Проверяем уникальность username (если изменился)
	if user.Username != entityUser.Username {
		var count int64
		r.db.WithContext(ctx).Model(&entity.User{}).
			Where("username = ? AND id != ?", user.Username, user.ID).
			Count(&count)
		if count > 0 {
			return ErrUsernameTaken
		}
	}

	return r.db.WithContext(ctx).Save(&entityUser).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var entityUsers []entity.User
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&entityUsers).Error
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(entityUsers))
	for i, eu := range entityUsers {
		users[i] = r.toDomain(&eu)
	}

	return users, nil
}

// Преобразование между доменной моделью и сущностью
func (r *userRepository) toEntity(domainUser *domain.User) *entity.User {
	return &entity.User{
		ID:           domainUser.ID,
		Email:        domainUser.Email,
		Username:     domainUser.Username,
		PasswordHash: domainUser.PasswordHash,
		Role:         domainUser.Role,
		CreatedAt:    domainUser.CreatedAt,
		UpdatedAt:    domainUser.UpdatedAt,
	}
}

func (r *userRepository) toDomain(entityUser *entity.User) *domain.User {
	return &domain.User{
		ID:           entityUser.ID,
		Email:        entityUser.Email,
		Username:     entityUser.Username,
		PasswordHash: entityUser.PasswordHash,
		Role:         entityUser.Role,
		CreatedAt:    entityUser.CreatedAt,
		UpdatedAt:    entityUser.UpdatedAt,
	}
}
