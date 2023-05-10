package storagememory

import (
	"github.com/google/uuid"
	"otus-microservices/hw06/internal/config"
	"otus-microservices/hw06/internal/storage"
	"sync"
)

type Storage struct {
	cfg      config.StorageConfig
	id       int64
	profiles map[int64]storage.Profile
	sessions map[string]storage.Session
	mutex    *sync.Mutex
}

func New(cfg config.StorageConfig) *Storage {
	mutex := sync.Mutex{}

	return &Storage{
		cfg:      cfg,
		id:       0,
		profiles: make(map[int64]storage.Profile),
		sessions: make(map[string]storage.Session),
		mutex:    &mutex,
	}
}

func (s *Storage) Connect() error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateProfile(profile storage.Profile) (int64, error) {
	s.mutex.Lock()
	s.id++
	profile.Id = s.id
	s.profiles[profile.Id] = profile
	s.mutex.Unlock()

	return profile.Id, nil
}

func (s *Storage) ReadProfile(id int64) (storage.Profile, error) {
	s.mutex.Lock()
	profile, ok := s.profiles[id]
	s.mutex.Unlock()

	if !ok {
		return storage.Profile{}, storage.ErrProfileNotFound
	}

	return profile, nil
}

func (s *Storage) UpdateProfile(id int64, profile storage.Profile) error {
	profile.Id = id
	s.mutex.Lock()
	s.profiles[profile.Id] = profile
	s.mutex.Unlock()

	return nil
}

func (s *Storage) DeleteProfile(id int64) error {
	_, err := s.ReadProfile(id)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	delete(s.profiles, id)
	s.mutex.Unlock()

	return nil
}

func (s *Storage) GetProfileByUsername(username string) (storage.Profile, error) {
	var (
		profileFound bool = false
		profile      storage.Profile
	)

	s.mutex.Lock()
	for _, value := range s.profiles {
		if value.Username == username {
			profileFound = true
			profile = value
			break
		}
	}
	s.mutex.Unlock()

	if !profileFound {
		return storage.Profile{}, storage.ErrProfileNotFound
	}

	return profile, nil
}

func (s *Storage) SetSessionForProfile(username string) (string, error) {
	profile, err := s.GetProfileByUsername(username)
	if err != nil {
		return "", err
	}

	sessionID := uuid.New().String()

	session := storage.Session{
		ProfileId: profile.Id,
		SessionId: sessionID,
	}

	s.mutex.Lock()
	s.sessions[sessionID] = session
	s.mutex.Unlock()

	return sessionID, nil
}

func (s *Storage) GetProfileForSession(sessionID string) (storage.Session, error) {
	s.mutex.Lock()
	session, ok := s.sessions[sessionID]
	s.mutex.Unlock()

	if !ok {
		return storage.Session{}, storage.ErrSessionNotFound
	}

	return session, nil
}

func (s *Storage) ClearSessionForProfile(username string) error {
	profile, err := s.GetProfileByUsername(username)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	for key, value := range s.sessions {
		if value.ProfileId == profile.Id {
			delete(s.sessions, key)
			break
		}
	}
	s.mutex.Unlock()

	return nil
}

func (s *Storage) ClearSessionForProfileId(id int64) error {
	s.mutex.Lock()
	for key, value := range s.sessions {
		if value.ProfileId == id {
			delete(s.sessions, key)
			break
		}
	}
	s.mutex.Unlock()

	return nil
}
