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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		s.config.Host, s.config.User, s.config.Password, s.config.Scheme, s.config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&storage.User{})

	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateUser(user storage.User) (int64, error) {
	result := s.db.Create(&user)

	return user.Id, result.Error
}

func (s *Storage) ReadUser(id int64) (storage.User, error) {
	user := storage.User{}
	result := s.db.First(&user, id)

	return user, result.Error
}

func (s *Storage) UpdateUser(id int64, user storage.User) error {
	user.Id = id
	result := s.db.Save(&user)
	return result.Error
}

func (s *Storage) DeleteUser(id int64) error {
	user := storage.User{Id: id}
	result := s.db.Delete(&user)
	return result.Error
}
