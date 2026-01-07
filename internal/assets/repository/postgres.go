package repository

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"invest-mate/internal/assets/models/entity"
	"invest-mate/internal/shared/config"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(cfg *config.Config) (*PostgresRepository, error) {
	dsn := cfg.GetDBDSN()

	log.Printf("Connecting to PostgreSQL (GORM): %s@%s:%d/%s",
		cfg.DBUser,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	if err := db.AutoMigrate(
		&entity.Bond{},
		&entity.Share{},
		&entity.Etf{},
		&entity.Currency{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	log.Println("✅ AutoMigrated (GORM)")

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) DB() *gorm.DB {
	return r.db
}

func SaveEntities[T entity.Marker](ctx context.Context, db *gorm.DB, entities []T) error {
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

func (r *PostgresRepository) GetBonds(ctx context.Context, limit, offset int) ([]entity.Bond, error) {
	var bonds []entity.Bond

	if err := r.db.WithContext(ctx).Order("ticker").Limit(limit).Offset(offset).Find(&bonds).Error; err != nil {
		return nil, fmt.Errorf("get bonds: %w", err)
	}

	return bonds, nil
}

func (r *PostgresRepository) GetShares(ctx context.Context, limit, offset int) ([]entity.Share, error) {
	var shares []entity.Share

	if err := r.db.WithContext(ctx).Order("ticker").Limit(limit).Offset(offset).Find(&shares).Error; err != nil {
		return nil, fmt.Errorf("get shares: %w", err)
	}

	return shares, nil
}

func (r *PostgresRepository) GetEtfs(ctx context.Context, limit, offset int) ([]entity.Etf, error) {
	var etfs []entity.Etf

	if err := r.db.WithContext(ctx).Order("ticker").Limit(limit).Offset(offset).Find(&etfs).Error; err != nil {
		return nil, fmt.Errorf("get etfs: %w", err)
	}

	return etfs, nil
}

func (r *PostgresRepository) GetCurrencies(ctx context.Context, limit, offset int) ([]entity.Currency, error) {
	var currencies []entity.Currency

	if err := r.db.WithContext(ctx).Order("ticker").Limit(limit).Offset(offset).Find(&currencies).Error; err != nil {
		return nil, fmt.Errorf("get currencies: %w", err)
	}

	return currencies, nil
}
