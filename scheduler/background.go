package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	UpdateDBTimeout  = 30 * time.Minute
	UpdateDBInterval = 6 * time.Hour
)

func UpdateDBScheduler(pool *pgxpool.Pool, logger *log.Logger) {
	ticker := time.NewTicker(UpdateDBInterval)
	defer ticker.Stop()
	ctx, close := context.WithTimeout(context.Background(), UpdateDBTimeout)
	defer close()

	for {
		<-ticker.C
		_, err := pool.Exec(ctx, "CALL refresh_mview()")
		if err != nil {
			logger.Println("Can't update database")
		}
	}
}
