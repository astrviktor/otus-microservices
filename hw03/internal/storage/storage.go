package storage

import "errors"

type User struct {
	Id        int64  `gorm:"id;uniqueIndex:idx_name" json:"id,omitempty"`
	Username  string `gorm:"user_name" json:"username,omitempty"`
	FirstName string `gorm:"first_name" json:"firstName,omitempty"`
	LastName  string `gorm:"last_name" json:"lastName,omitempty"`
	Email     string `gorm:"email" json:"email,omitempty"`
	Phone     string `gorm:"phone" json:"phone,omitempty"`
}

func (u *User) Validate() error {
	return nil
}

type Storage interface {
	Connect() error
	Close() error

	CreateUser(user User) (int64, error)
	ReadUser(id int64) (User, error)
	UpdateUser(id int64, user User) error
	DeleteUser(id int64) error
}

var (
	ErrUserNotFound = errors.New("user not found in storage\n")
)
