package storagegorm

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"otus-microservices/billing/internal/config"
	"otus-microservices/billing/internal/storage"
)

type Storage struct {
	db     *gorm.DB
	config config.StorageConfig
}

func New(cfg config.StorageConfig) *Storage {
	return &Storage{
		config: cfg,
	}
}

func (s *Storage) Connect() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		s.config.Host, s.config.User, s.config.Password, s.config.Scheme, s.config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&storage.Billing{})

	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateBilling(billing storage.Billing) (int64, error) {
	getBilling := storage.Billing{}
	result := s.db.Where("client_id = ?", billing.ClientID).First(&getBilling)

	if result.Error != nil {
		result = s.db.Create(&billing)
		return billing.ClientID, result.Error
	}

	return getBilling.ClientID, result.Error
}

func (s *Storage) ReadBilling(id int64) (storage.Billing, error) {
	billing := storage.Billing{}
	result := s.db.Where("client_id = ?", id).First(&billing)

	return billing, result.Error
}

func (s *Storage) UpdateBilling(id int64, billing storage.Billing) error {
	billing.ClientID = id
	result := s.db.Model(&billing).Where("client_id = ?", id).
		Update("balance", billing.Balance)
	return result.Error
}

func (s *Storage) DeleteBilling(id int64) error {
	billing := storage.Billing{ClientID: id}
	result := s.db.Delete(&billing)
	return result.Error
}
