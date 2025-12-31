package models

import "github.com/google/uuid"

type FavoriteMovie struct {
	UserID  uuid.UUID `json:"user_id"`
	MovieID uuid.UUID `json:"movie_id"`
}

func (fm FavoriteMovie) ToResponse() FavoriteMovieResponse {
	return FavoriteMovieResponse{
		MovieID: fm.MovieID,
	}
}

type FavoriteMovieList struct {
	FavMovieList []FavoriteMovie `json:"fav_movie_list"`
}

func (ml *FavoriteMovieList) ToResponse() FavoriteMovieListResponse {
	favMovieListResponse := []FavoriteMovieResponse{}

	for _, favMovie := range ml.FavMovieList {
		favMovieListResponse = append(favMovieListResponse, favMovie.ToResponse())
	}
	return FavoriteMovieListResponse{
		FavMovieList: favMovieListResponse,
	}
}

type FavoriteMovieResponse struct {
	MovieID uuid.UUID `json:"movie_id"`
}

type FavoriteMovieListResponse struct {
	FavMovieList []FavoriteMovieResponse `json:"movie_list"`
}
