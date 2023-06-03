package storage

import (
	"errors"
)

type Order struct {
	Id        int64   `gorm:"id;uniqueIndex:idx_name" json:"id,omitempty"`
	RequestId string  `gorm:"request_id;not null;unique" json:"-"`
	Name      string  `gorm:"name" json:"name,omitempty"`
	Email     string  `gorm:"email" json:"email,omitempty"`
	Total     float32 `gorm:"total" json:"total,omitempty"`
}

func (u *Order) Validate() error {
	return nil
}

type Storage interface {
	Connect() error
	Close() error

	CreateOrder(user Order) (int64, error)
	ReadOrder(id int64) (Order, error)
	UpdateOrder(id int64, user Order) error
	DeleteOrder(id int64) error
}

var (
	ErrOrderNotFound = errors.New("order not found in storage\n")
)
