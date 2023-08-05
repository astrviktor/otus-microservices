package storagegorm

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"otus-microservices/orderservice/internal/config"
	"otus-microservices/orderservice/internal/storage"
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

func (s *Storage) CreateOrder(order storage.Order) (int64, error) {
	result := s.db.Create(&order)
	return order.Id, result.Error
}

func (s *Storage) ReadOrder(id int64) (storage.Order, error) {
	order := storage.Order{}
	result := s.db.First(&order, id)

	return order, result.Error
}

func (s *Storage) UpdateOrder(id int64, order storage.Order) error {
	order.Id = id
	result := s.db.Save(&order)
	return result.Error
}

func (s *Storage) DeleteOrder(id int64) error {
	order := storage.Order{Id: id}
	result := s.db.Delete(&order)
	return result.Error
}
