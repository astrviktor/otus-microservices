package storage

import (
	"errors"
)

type Notification struct {
	ClientID int64  `gorm:"client_id" json:"client_id"`
	OrderID  int64  `gorm:"order_id" json:"order_id"`
	Theme    string `gorm:"theme" json:"theme"`
	Message  string `gorm:"message" json:"message"`
}

func (u *Notification) Validate() error {
	return nil
}

type Storage interface {
	Connect() error
	Close() error

	CreateNotification(notification Notification) error
	ReadNotification(orderID int64) (Notification, error)
}

var (
	ErrNotificationNotFound = errors.New("notification not found in storage\n")
)
