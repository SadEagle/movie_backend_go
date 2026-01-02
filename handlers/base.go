package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func writeResponseBody[T any](rw http.ResponseWriter, responseByteObj T, responseObjName string) {
	responseFavMovieByte, err := json.Marshal(responseByteObj)
	if err != nil {
		log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't convert %s data to response json", responseObjName), 500)
		return
	}

	_, err = rw.Write(responseFavMovieByte)
	if err != nil {
		log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't send %s data", responseObjName), 500)
		return
	}
}
