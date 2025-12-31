package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ResponseConvertable[ResponseT any] interface {
	ToResponse() ResponseT
}

func writeResponseBody[ResponseT any](rw http.ResponseWriter, db_object ResponseConvertable[ResponseT], responseObjName string) {
	responseFavMovie := db_object.ToResponse()
	responseFavMovieByte, err := json.Marshal(responseFavMovie)
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
