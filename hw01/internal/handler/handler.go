package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ResponseHealth struct {
	Status string `json:"status"`
}

func WriteResponse(w http.ResponseWriter, resp interface{}) {
	respBuf, err := json.Marshal(resp)
	if err != nil {
		log.Println(fmt.Sprintf("response marshal error: %s", err))
	}

	respBuf = append(respBuf, []byte("\n")...)
	_, err = w.Write(respBuf)

	if err != nil {
		log.Println(fmt.Sprintf("response write error: %s", err))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		w.WriteHeader(http.StatusOK)
		WriteResponse(w, &ResponseHealth{
			Status: "OK",
		})
	}

	return
}
