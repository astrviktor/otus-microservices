package storage

import (
	"errors"
)

type Billing struct {
	ClientID int64 `gorm:"client_id;uniqueIndex:idx_name" json:"client_id,omitempty"`
	Balance  int64 `gorm:"balance" json:"balance"`
}

func (u *Billing) Validate() error {
	return nil
}

type Storage interface {
	Connect() error
	Close() error

	CreateBilling(billing Billing) (int64, error)
	ReadBilling(id int64) (Billing, error)
	UpdateBilling(id int64, user Billing) error
	DeleteBilling(id int64) error
}

var (
	ErrBillingNotFound = errors.New("order not found in storage\n")
)
