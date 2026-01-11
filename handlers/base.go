package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HandlerObj struct {
	DB  *sql.DB
	Log log.Logger
}

func writeResponseBody(rw http.ResponseWriter, responseByteObj any, responseObjName string) {
	var enc = json.NewEncoder(rw)
	err := enc.Encode(responseByteObj)
	if err != nil {
		log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't send %s data", responseObjName), 500)
		return
	}
}
