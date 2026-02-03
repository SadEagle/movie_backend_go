package reqmodel

import "github.com/jackc/pgx/v5/pgtype"

type MovieFavoriteListResponse struct {
	MovieID         pgtype.UUID   `json:"movie_id"`
	FavoriteUserIDs []pgtype.UUID `json:"favorite_movie_ids"`
}

type UserFavoriteListResponse struct {
	UserID           pgtype.UUID   `json:"user_id"`
	FavoriteMovieIDs []pgtype.UUID `json:"favorite_user_ids"`
}

type FavoriteGetRequest struct {
	UserID  pgtype.UUID `json:"user_id"`
	MovieID pgtype.UUID `json:"movie_id"`
}

type FavoriteCreateRequest struct {
	MovieID pgtype.UUID `json:"movie_id"`
}

type FavoriteDeleteRequest struct {
	MovieID pgtype.UUID `json:"movie_id"`
}
