package storagegorm

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"otus-microservices/notification/internal/config"
	"otus-microservices/notification/internal/storage"
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

	err = db.AutoMigrate(&storage.Notification{})

	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateNotification(Notification storage.Notification) error {
	getNotification := storage.Notification{}
	result := s.db.Where("client_id = ? and order_id = ?",
		Notification.ClientID, Notification.OrderID).First(&getNotification)

	if result.Error != nil {
		result = s.db.Create(&Notification)
		return result.Error
	}

	return result.Error
}

func (s *Storage) ReadNotification(orderID int64) (storage.Notification, error) {
	notification := storage.Notification{}
	result := s.db.Where("order_id = ?", orderID).First(&notification)

	return notification, result.Error
}
