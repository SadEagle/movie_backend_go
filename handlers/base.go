package handlers

import (
	"movie_backend_go/db/sqlc"
	"time"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const OpTimeContext = 5 * time.Minute
const CheckHealthTimeContext = 2 * time.Minute

type HandlerObj struct {
	DBPool sqlc.Querier
	Log    *log.Logger
}

func writeResponseBody(rw http.ResponseWriter, responseByteObj any, responseObjName string) {
	rw.Header().Set("Content-Type", "application/json")
	var enc = json.NewEncoder(rw)
	err := enc.Encode(responseByteObj)
	if err != nil {
		log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't send %s data", responseObjName), http.StatusInternalServerError)
		return
	}
}
