package crud

import (
	"errors"
	"github.com/google/uuid"
	"movie_backend_go/model"
	// "reflect"
	"sync"
)

type FavoriteMovieStorageType struct {
	favoriteMovie []model.FavoriteMovie
	mu            sync.RWMutex
}

func (m *FavoriteMovieStorageType) Add(favoriteMovieData model.FavoriteMovie) (model.FavoriteMovie, error) {
	m.mu.Lock()
	m.favoriteMovie = append(m.favoriteMovie, favoriteMovieData)
	m.mu.Unlock()
	return favoriteMovieData, nil
}

func (m *FavoriteMovieStorageType) getElementByID(userID uuid.UUID) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var favMovIdx int
	for favMovIdx = range len(m.favoriteMovie) {
		if m.favoriteMovie[favMovIdx].ID == userID {
			return favMovIdx, nil
		}
	}
	return -1, errors.New("Data wasn't found")
}

func (m *FavoriteMovieStorageType) ReadUser(id uuid.UUID) (model.FavoriteMovie, error) {
	favMovIdx, error := m.getElementByID(id)
	if error != nil {
		return model.FavoriteMovie{}, error
	}
	return m.favoriteMovie[favMovIdx], nil
}

func (m *FavoriteMovieStorageType) ReadMovie(id uuid.UUID) (model.FavoriteMovie, error) {
	favMovIdx, error := m.getElementByID(id)
	if error != nil {
		return model.FavoriteMovie{}, error
	}
	return m.favoriteMovie[favMovIdx], nil
}

func (m *FavoriteMovieStorageType) Delete(id uuid.UUID) error {
	user_idx, error := m.getElementByID(id)
	if error != nil {
		return error
	}

	m.mu.Lock()
	if user_idx == 0 {
		m.favoriteMovie = m.favoriteMovie[1:]
	} else if len(m.favoriteMovie) == user_idx-1 {
		m.favoriteMovie = m.favoriteMovie[:user_idx]
	} else {
		m.favoriteMovie = append(m.favoriteMovie[:user_idx], m.favoriteMovie[user_idx+1:]...)
	}
	m.mu.Unlock()
	return nil
}
