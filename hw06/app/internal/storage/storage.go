package storage

import (
	"errors"
)

type Profile struct {
	Id       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Session struct {
	ProfileId int64
	SessionId string
}

type Storage interface {
	Connect() error
	Close() error

	CreateProfile(profile Profile) (int64, error)
	ReadProfile(id int64) (Profile, error)
	UpdateProfile(id int64, profile Profile) error
	DeleteProfile(id int64) error

	//CheckProfile(profile Profile) (int64, error)

	SetSessionForProfile(username string) (string, error)
	GetProfileForSession(sessionID string) (Session, error)
	ClearSessionForProfile(username string) error
	ClearSessionForProfileId(id int64) error
}

var (
	ErrProfileNotFound = errors.New("profile not found\n")
	ErrSessionNotFound = errors.New("session not found\n")
)
