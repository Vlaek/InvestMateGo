package repository

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	config "invest-mate/config"
	repoModels "invest-mate/internal/repository/models"
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
		&repoModels.Bond{},
		&repoModels.Share{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	log.Println("✅ AutoMigrated (GORM)")

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) DB() *gorm.DB {
	return r.db
}

func SaveEntities[T repoModels.RepositoryMarker](ctx context.Context, db *gorm.DB, entities []T) error {
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

// GetBonds возвращает облигации из БД
func (r *PostgresRepository) GetBonds(ctx context.Context, limit, offset int) ([]repoModels.Bond, error) {
	var bonds []repoModels.Bond
	if err := r.db.WithContext(ctx).Order("ticker").Limit(limit).Offset(offset).Find(&bonds).Error; err != nil {
		return nil, fmt.Errorf("get bonds: %w", err)
	}
	return bonds, nil
}

// GetBondByTicker ищет облигацию по тикеру
func (r *PostgresRepository) GetBondByTicker(ctx context.Context, ticker string) (*repoModels.Bond, error) {
	var bond repoModels.Bond
	err := r.db.WithContext(ctx).Where("ticker = ?", ticker).First(&bond).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get bond by ticker: %w", err)
	}
	return &bond, nil
}

// GetBondCount возвращает количество облигаций
func (r *PostgresRepository) GetBondCount(ctx context.Context) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&repoModels.Bond{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count bonds: %w", err)
	}
	return int(count), nil
}
