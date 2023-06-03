package storagegorm

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"otus-microservices/hw03/internal/config"
	"otus-microservices/hw03/internal/storage"
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

	err = db.AutoMigrate(&storage.Order{})

	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateOrder(Order storage.Order) (int64, error) {
	GetOrder := storage.Order{}
	result := s.db.Where("request_id = ?", Order.RequestId).First(&GetOrder)

	if result.Error != nil {
		result = s.db.Create(&Order)
		return Order.Id, result.Error
	}

	return GetOrder.Id, result.Error
}

func (s *Storage) ReadOrder(id int64) (storage.Order, error) {
	Order := storage.Order{}
	result := s.db.First(&Order, id)

	return Order, result.Error
}

func (s *Storage) UpdateOrder(id int64, Order storage.Order) error {
	Order.Id = id
	result := s.db.Save(&Order)
	return result.Error
}

func (s *Storage) DeleteOrder(id int64) error {
	Order := storage.Order{Id: id}
	result := s.db.Delete(&Order)
	return result.Error
}
