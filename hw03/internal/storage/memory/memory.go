package storagememory

import (
	"otus-microservices/hw03/internal/config"
	"otus-microservices/hw03/internal/storage"
	"sync"
)

type Storage struct {
	cfg   config.StorageConfig
	id    int64
	users map[int64]storage.User
	mutex *sync.Mutex
}

func New(cfg config.StorageConfig) *Storage {
	mutex := sync.Mutex{}

	return &Storage{
		cfg:   cfg,
		id:    0,
		users: make(map[int64]storage.User),
		mutex: &mutex,
	}
}

func (s *Storage) Connect() error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateUser(user storage.User) (int64, error) {
	s.mutex.Lock()
	s.id++
	user.Id = s.id
	s.users[user.Id] = user
	s.mutex.Unlock()

	return user.Id, nil
}

func (s *Storage) ReadUser(id int64) (storage.User, error) {
	s.mutex.Lock()
	user, ok := s.users[id]
	s.mutex.Unlock()

	if !ok {
		return storage.User{}, storage.ErrUserNotFound
	}

	return user, nil
}

func (s *Storage) UpdateUser(id int64, user storage.User) error {
	user.Id = id
	s.mutex.Lock()
	s.users[user.Id] = user
	s.mutex.Unlock()

	return nil
}

func (s *Storage) DeleteUser(id int64) error {
	_, err := s.ReadUser(id)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	delete(s.users, id)
	s.mutex.Unlock()

	return nil
}
