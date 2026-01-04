package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func writeResponseBody[T any](rw http.ResponseWriter, responseByteObj T, responseObjName string) {
	var enc = json.NewEncoder(rw)
	err := enc.Encode(responseByteObj)
	if err != nil {
		log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't send %s data", responseObjName), 500)
		return
	}
}
