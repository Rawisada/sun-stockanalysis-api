package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type StockRepository interface {
	FindByID(id uuid.UUID) (*models.Stock, error)
	Create(stock *models.Stock) error
	ListSymbols() ([]string, error)
	FindAll() ([]models.Stock, error)
	EnsureMasterAssetType(name string) error
	EnsureMasterExchange(name string) error
	EnsureMasterSector(name string) error
}

type StockRepositoryImpl struct {
	db *gorm.DB // GORM DB instance
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &StockRepositoryImpl{db: db}
}

func (r *StockRepositoryImpl) FindByID(id uuid.UUID) (*models.Stock, error) {
	var s models.Stock
	if err := r.db.First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StockRepositoryImpl) Create(s *models.Stock) error {
	return r.db.Create(s).Error
}

func (r *StockRepositoryImpl) ListSymbols() ([]string, error) {
	var symbols []string
	if err := r.db.Model(&models.Stock{}).
		Where("is_active = ?", true).
		Pluck("symbol", &symbols).Error; err != nil {
		return nil, err
	}
	return symbols, nil
}

func (r *StockRepositoryImpl) FindAll() ([]models.Stock, error) {
	var stocks []models.Stock
	if err := r.db.
		Order("created_at desc").
		Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *StockRepositoryImpl) EnsureMasterAssetType(name string) error {
	if name == "" {
		return nil
	}
	record := models.MasterAssetType{Name: name, IsActive: true}
	return r.db.Where("name = ?", name).FirstOrCreate(&record).Error
}

func (r *StockRepositoryImpl) EnsureMasterExchange(name string) error {
	if name == "" {
		return nil
	}
	record := models.MasterExchange{Name: name, IsActive: true}
	return r.db.Where("name = ?", name).FirstOrCreate(&record).Error
}

func (r *StockRepositoryImpl) EnsureMasterSector(name string) error {
	if name == "" {
		return nil
	}
	record := models.MasterSector{Name: name, IsActive: true}
	return r.db.Where("name = ?", name).FirstOrCreate(&record).Error
}
