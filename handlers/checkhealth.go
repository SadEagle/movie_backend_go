package handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// @Summary      Dummy healthcheck
// @Description  Healtcheck dummy output
// @Tags         healthcheck
// @Accept       json
// @Produce      json
// @Success      200  {object}  sqlc.GetMovieByIDRow
// @Failure      503  {object}  map[string]string
// @Router       /healthcheck [get]
func CheckHealthHandlerCreate(pool *pgxpool.Pool) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), CheckHealthTimeContext)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			http.Error(rw, "Service connection lost", http.StatusServiceUnavailable)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
