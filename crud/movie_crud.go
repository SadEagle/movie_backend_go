package crud

import (
	"errors"
	"github.com/google/uuid"
	"movie_backend_go/model"
	// "reflect"
	"sync"
)

type MovieStorageType struct {
	movie []model.Movie
	mu    sync.RWMutex
}

func (m *MovieStorageType) Add(movieData model.Movie) (model.Movie, error) {
	m.mu.Lock()
	m.movie = append(m.movie, movieData)
	m.mu.Unlock()
	return movieData, nil
}

func (m *MovieStorageType) getElementByID(id uuid.UUID) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var movieIdx int
	for movieIdx = range len(m.movie) {
		if m.movie[movieIdx].ID == id {
			return movieIdx, nil
		}
	}
	return -1, errors.New("Data wasn't found")
}

func (m *MovieStorageType) Read(id uuid.UUID) (model.Movie, error) {
	movieIdx, error := m.getElementByID(id)
	if error != nil {
		return model.Movie{}, error
	}
	return m.movie[movieIdx], nil
}

func (m *MovieStorageType) ReadList() ([]model.Movie, error) {
	return m.movie, nil
}

func (m *MovieStorageType) Update(id uuid.UUID, updateData model.MovieUpdateRequest) (model.Movie, error) {
	movieIdx, error := m.getElementByID(id)
	if error != nil {
		return model.Movie{}, error
	}

	// TODO: add generic with reflect instead of dummy update
	if updateData.Title != nil {
		m.movie[movieIdx].Title = *updateData.Title
	}
	if updateData.Rating != nil {
		m.movie[movieIdx].Rating = *updateData.Rating
	}
	return m.movie[movieIdx], nil
}

func (m *MovieStorageType) Delete(id uuid.UUID) error {
	movieIdx, error := m.getElementByID(id)
	if error != nil {
		return error
	}

	m.mu.Lock()
	if movieIdx == 0 {
		m.movie = m.movie[1:]
	} else if len(m.movie) == movieIdx-1 {
		m.movie = m.movie[:movieIdx]
	} else {
		m.movie = append(m.movie[:movieIdx], m.movie[movieIdx+1:]...)
	}
	m.mu.Unlock()
	return nil
}
