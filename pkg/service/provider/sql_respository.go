package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
)

// ProviderSQLModel is a SQL database model for provider.
type ProviderSQLModel struct {
	ID        uint       `gorm:"column:id;PRIMARY_KEY"`
	UUID      string     `gorm:"column:uuid;UNIQUE;UNIQUE_INDEX;NOT NULL"`
	ShortName string     `gorm:"column:short_name;INDEX;NOT NULL"`
	LongName  string     `gorm:"column:long_name;NOT NULL"`
	CreatedAt time.Time  `gorm:"column:created_at;NOT NULL"`
	UpdatedAt time.Time  `gorm:"column:updated_at;NOT NULL"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

// TableName sets provider table name.
func (pm *ProviderSQLModel) TableName() string {
	return "provider"
}

func (pm *ProviderSQLModel) toProvider() *domain.Provider {
	return &domain.Provider{
		UUID:      pm.UUID,
		ShortName: pm.ShortName,
		LongName:  pm.LongName,
	}
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates new provider repository.
func NewRepository(db *gorm.DB, automigrate bool) domain.ProviderRepository {
	if automigrate {
		db.AutoMigrate(&ProviderSQLModel{})
	}
	return &repository{db}
}

// CreateProvider creates new provider in the database.
func (r *repository) CreateProvider(ctx context.Context, p domain.Provider) error {
	if err := r.db.Create(&ProviderSQLModel{
		UUID:      p.UUID,
		ShortName: p.ShortName,
		LongName:  p.LongName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}).Error; err != nil {
		log.Error(err, "Failed creating new provider")
		return err
	}
	return nil
}

// UpdateProvider updates the existing provider in the database.
func (r *repository) UpdateProvider(ctx context.Context, p domain.Provider) error {
	if err := r.db.Model(&ProviderSQLModel{}).
		Where("uuid = ? AND deleted_at IS NULL", p.UUID).Updates(
		map[string]interface{}{
			"short_name": p.ShortName,
			"long_name":  p.LongName,
			"updated_at": time.Now(),
		},
	).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Warn(fmt.Sprintf("Provider with '%s' UUID doesn't exist", p.UUID))
			return domain.ErrNotFound
		}
		log.Error(err, fmt.Sprintf("Failed updating provider with '%s' UUID", p.UUID))
		return err
	}
	return nil
}

// DeleteProviderByUUID deletes existing provider in the database based on its UUID.
func (r *repository) DeleteProviderByUUID(ctx context.Context, uuid string) error {
	if err := r.db.Where("uuid = ?", uuid).Delete(&ProviderSQLModel{}).Error; err != nil {
		log.Error(err, fmt.Sprintf("Failed deleting provider with '%s' UUID", uuid))
		return err
	}
	return nil
}

// GetProviderByUUID gets a provider in the database based on its UUID.
func (r *repository) GetProviderByUUID(ctx context.Context, uuid string) (*domain.Provider, error) {
	var pm ProviderSQLModel
	if err := r.db.Where("uuid = ? AND deleted_at IS NULL", uuid).First(&pm).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Warn(fmt.Sprintf("Provider with '%s' UUID doesn't exist", uuid))
			return nil, domain.ErrNotFound
		}
		log.Error(err, fmt.Sprintf("Failed getting provider with '%s' UUID", uuid))
		return nil, err
	}
	return pm.toProvider(), nil
}

// GetProviderByShortName gets a provider in the database based on its short name.
func (r *repository) GetProviderByShortName(ctx context.Context, shortName string) (*domain.Provider, error) {
	var pm ProviderSQLModel
	if err := r.db.Where("short_name = ? AND deleted_at IS NULL", shortName).First(&pm).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Warn(fmt.Sprintf("Provider with '%s' short name doesn't exist", shortName))
			return nil, domain.ErrNotFound
		}
		log.Error(err, fmt.Sprintf("Failed getting provider with '%s' short name", shortName))
		return nil, err
	}
	return pm.toProvider(), nil
}

// GetProviders gets all providers in the database.
func (r *repository) GetProviders(ctx context.Context, limit int) ([]domain.Provider, error) {
	var pms []ProviderSQLModel
	query := r.db
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&pms).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		log.Error(err, "Failed getting providers")
		return nil, err
	}
	results := []domain.Provider{}
	for _, pm := range pms {
		results = append(results, domain.Provider{
			UUID:      pm.UUID,
			ShortName: pm.ShortName,
			LongName:  pm.LongName,
		})
	}
	return results, nil
}