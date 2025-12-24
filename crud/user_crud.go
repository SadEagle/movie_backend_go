package crud

import (
	"errors"
	"github.com/google/uuid"
	"movie_backend_go/model"
	// "reflect"
	"sync"
)

type UserStorageType struct {
	user []model.User
	mu   sync.RWMutex
}

func (m *UserStorageType) Add(userData model.User) (model.User, error) {
	m.mu.Lock()
	m.user = append(m.user, userData)
	m.mu.Unlock()
	return userData, nil
}

func (m *UserStorageType) getElementByID(id uuid.UUID) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var userIdx int
	for userIdx = range len(m.user) {
		if m.user[userIdx].ID == id {
			return userIdx, nil
		}
	}
	return -1, errors.New("Data wasn't found")
}

func (m *UserStorageType) Read(id uuid.UUID) (model.User, error) {
	userIdx, error := m.getElementByID(id)
	if error != nil {
		return model.User{}, error
	}
	return m.user[userIdx], nil
}

func (m *UserStorageType) ReadList() ([]model.User, error) {
	return m.user, nil
}

func (m *UserStorageType) Update(id uuid.UUID, updateData model.UserUpdateRequest) (model.User, error) {
	userIdx, error := m.getElementByID(id)
	if error != nil {
		return model.User{}, error
	}

	// TODO: add generic with reflect instead of dummy update
	if updateData.Login != nil {
		m.user[userIdx].Login = *updateData.Login
	}
	if updateData.Name != nil {
		m.user[userIdx].Name = *updateData.Name
	}
	if updateData.Password != nil {
		m.user[userIdx].Password = *updateData.Password
	}
	return m.user[userIdx], nil
}

func (m *UserStorageType) Delete(id uuid.UUID) error {
	userIdx, error := m.getElementByID(id)
	if error != nil {
		return error
	}

	m.mu.Lock()
	if userIdx == 0 {
		m.user = m.user[1:]
	} else if len(m.user) == userIdx-1 {
		m.user = m.user[:userIdx]
	} else {
		m.user = append(m.user[:userIdx], m.user[userIdx+1:]...)
	}
	m.mu.Unlock()
	return nil
}
