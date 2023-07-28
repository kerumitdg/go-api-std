// these are all dummy store functions for testing purposes

package stores

import (
	"errors"

	"github.com/fredrikaverpil/go-api-std/lib"
	"github.com/fredrikaverpil/go-api-std/models"
)

type DummyDbRecord struct {
	ID             int
	Username       string
	HashedPassword string
}

type DummyStore struct {
	userDb map[int]DummyDbRecord
}

func NewDummyStore() DummyStore {
	userDb := make(map[int]DummyDbRecord)
	return DummyStore{userDb: userDb}
}

func (s *DummyStore) CreateUser(username string, password string) (models.User, error) {
	nextID := len(s.userDb) + 1
	hashedPassword, err := lib.HashPassword(password)
	if err != nil {
		return models.User{}, errors.New("could not hash password")
	}

	user := models.User{ID: nextID, Username: username}
	err = user.Validate()
	if err != nil {
		return models.User{}, errors.New("could not validate user")
	}

	userRecord := DummyDbRecord{ID: nextID, Username: username, HashedPassword: hashedPassword}
	s.userDb[nextID] = userRecord

	return user, nil
}

func (s *DummyStore) GetUserByUsername(username string) (models.User, error) {
	for userId, dbRecord := range s.userDb {
		if dbRecord.Username == username {
			user, err := s.GetUser(userId)
			if err != nil {
				return models.User{}, errors.New("could not get user")
			}
			return user, nil
		}
	}

	return models.User{}, errors.New("no user found")
}

func (s *DummyStore) GetUser(id int) (models.User, error) {
	userRecord, exists := s.userDb[id]
	if !exists {
		return models.User{}, errors.New("no user found")
	}

	user := models.User{ID: userRecord.ID, Username: userRecord.Username}
	err := user.Validate()
	if err != nil {
		return models.User{}, errors.New("could not validate user")
	}

	return user, nil
}
