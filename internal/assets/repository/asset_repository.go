package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"invest-mate/internal/assets/models/entity"
	"invest-mate/pkg/logger"
)

type AssetRepository interface {
	GetDB() *gorm.DB

	GetAssetByField(ctx context.Context, fieldName string, fieldValue string) (*entity.Asset, error)

	GetBonds(ctx context.Context, limit, offset int) ([]entity.Bond, error)
	GetBondByField(ctx context.Context, fieldName string, fieldValue string) (*entity.Bond, error)

	GetShares(ctx context.Context, limit, offset int) ([]entity.Share, error)
	GetShareByField(ctx context.Context, fieldName string, fieldValue string) (*entity.Share, error)

	GetEtfs(ctx context.Context, limit, offset int) ([]entity.Etf, error)
	GetEtfByField(ctx context.Context, fieldName string, fieldValue string) (*entity.Etf, error)

	GetCurrencies(ctx context.Context, limit, offset int) ([]entity.Currency, error)
	GetCurrencyByField(ctx context.Context, fieldName string, fieldValue string) (*entity.Currency, error)
}

type assetRepository struct {
	db *gorm.DB
}

// Создание нового репозитория
func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

// Получение БД
func (r *assetRepository) GetDB() *gorm.DB {
	return r.db
}

// Сохранение сущностей в БД
func SaveToDB[T entity.Marker](ctx context.Context, db *gorm.DB, entities []T, entityName string) error {
	if len(entities) == 0 {
		return nil
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.ErrorLog("Panic during %s save: %v", entityName, r)
		}
	}()

	if err := saveEntities(ctx, tx, entities); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save %s: %w", entityName, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit %s transaction: %w", entityName, err)
	}

	logger.InfoLog("Saved %d %s to database", len(entities), entityName)
	return nil
}

// Сохранение сущностей в БД
func saveEntities[T entity.Marker](ctx context.Context, db *gorm.DB, entities []T) error {
	if len(entities) == 0 {
		return nil
	}

	const batchSize = 1000

	for i := 0; i < len(entities); i += batchSize {
		end := i + batchSize

		if end > len(entities) {
			end = len(entities)
		}

		if err := db.WithContext(ctx).Create(entities[i:end]).Error; err != nil {
			return err
		}
	}

	return nil
}

// Получение инструмента по полю из БД
func (r *assetRepository) GetAssetByField(ctx context.Context, fieldName string, fieldValue string) (*entity.Asset, error) {
	var entity entity.Asset
	query := fmt.Sprintf("%s = ?", fieldName)
	result := r.db.Where(query, fieldValue).First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &entity, nil
}

// Получение облигаций из БД
func (r *assetRepository) GetBonds(ctx context.Context, limit int, offset int) ([]entity.Bond, error) {
	var bonds []entity.Bond

	if err := r.db.WithContext(ctx).Order("ticker").Limit(limit).Offset(offset).Find(&bonds).Error; err != nil {
		return nil, fmt.Errorf("get bonds: %w", err)
	}

	return bonds, nil
}

// Получение облигации по полю из БД
func (r *assetRepository) GetBondByField(ctx context.Context, fieldName string, fieldValue string) (*entity.Bond, error) {
	var entity entity.Bond
	query := fmt.Sprintf("%s = ?", fieldName)
	result := r.db.Where(query, fieldValue).First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &entity, nil
}

// Получение акций из БД
func (r *assetRepository) GetShares(ctx context.Context, limit int, offset int) ([]entity.Share, error) {
	var shares []entity.Share

	if err := r.db.WithContext(ctx).Order("ticker").Limit(limit).Offset(offset).Find(&shares).Error; err != nil {
		return nil, fmt.Errorf("get shares: %w", err)
	}

	return shares, nil
}

// Получение акции по полю из БД
func (r *assetRepository) GetShareByField(ctx context.Context, fieldName string, fieldValue string) (*entity.Share, error) {
	var entity entity.Share
	query := fmt.Sprintf("%s = ?", fieldName)
	result := r.db.Where(query, fieldValue).First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &entity, nil
}

// Получение фондов из БД
func (r *assetRepository) GetEtfs(ctx context.Context, limit int, offset int) ([]entity.Etf, error) {
	var etfs []entity.Etf

	if err := r.db.WithContext(ctx).Order("ticker").Limit(limit).Offset(offset).Find(&etfs).Error; err != nil {
		return nil, fmt.Errorf("get etfs: %w", err)
	}

	return etfs, nil
}

// Получение фонда по идентификатору из БД
func (r *assetRepository) GetEtfByField(ctx context.Context, fieldName string, fieldValue string) (*entity.Etf, error) {
	var entity entity.Etf
	query := fmt.Sprintf("%s = ?", fieldName)
	result := r.db.Where(query, fieldValue).First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &entity, nil
}

// Получение валют из БД
func (r *assetRepository) GetCurrencies(ctx context.Context, limit int, offset int) ([]entity.Currency, error) {
	var currencies []entity.Currency

	if err := r.db.WithContext(ctx).Order("ticker").Limit(limit).Offset(offset).Find(&currencies).Error; err != nil {
		return nil, fmt.Errorf("get currencies: %w", err)
	}

	return currencies, nil
}

// Получение валюты по идентификатору из БД
func (r *assetRepository) GetCurrencyByField(ctx context.Context, fieldName string, fieldValue string) (*entity.Currency, error) {
	var entity entity.Currency
	query := fmt.Sprintf("%s = ?", fieldName)
	result := r.db.Where(query, fieldValue).First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &entity, nil
}
