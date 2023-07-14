package storagememory

import (
	"otus-microservices/orderservice/internal/config"
	"otus-microservices/orderservice/internal/storage"
	"sync"
)

type Storage struct {
	cfg    config.StorageConfig
	id     int64
	Orders map[int64]storage.Order
	mutex  *sync.Mutex
}

func New(cfg config.StorageConfig) *Storage {
	mutex := sync.Mutex{}

	return &Storage{
		cfg:    cfg,
		id:     0,
		Orders: make(map[int64]storage.Order),
		mutex:  &mutex,
	}
}

func (s *Storage) Connect() error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateOrder(Order storage.Order) (int64, error) {
	s.mutex.Lock()
	s.id++
	Order.Id = s.id
	s.Orders[Order.Id] = Order
	s.mutex.Unlock()

	return Order.Id, nil
}

func (s *Storage) ReadOrder(id int64) (storage.Order, error) {
	s.mutex.Lock()
	Order, ok := s.Orders[id]
	s.mutex.Unlock()

	if !ok {
		return storage.Order{}, storage.ErrOrderNotFound
	}

	return Order, nil
}

func (s *Storage) UpdateOrder(id int64, Order storage.Order) error {
	Order.Id = id
	s.mutex.Lock()
	s.Orders[Order.Id] = Order
	s.mutex.Unlock()

	return nil
}

func (s *Storage) DeleteOrder(id int64) error {
	_, err := s.ReadOrder(id)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	delete(s.Orders, id)
	s.mutex.Unlock()

	return nil
}
