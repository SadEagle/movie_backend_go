package handlers

import (
	"movie_backend_go/db/sqlc"
	"time"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	OpTimeContext          = 5 * time.Minute
	CheckHealthTimeContext = 2 * time.Minute
)

type HandlerObj struct {
	QuerierDB *sqlc.Queries
	Logger    *log.Logger
}

func writeResponseBody(rw http.ResponseWriter, responseObj any, responseObjName string) {
	rw.Header().Set("Content-Type", "application/json")
	var enc = json.NewEncoder(rw)
	err := enc.Encode(responseObj)
	if err != nil {
		log.Println(err)
		http.Error(rw, fmt.Sprintf("can't send %s data", responseObjName), http.StatusInternalServerError)
		return
	}
}
