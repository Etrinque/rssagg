package utils

// Json respose for endpoint
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type JsonUtils struct {
	RespondWithError http.HandlerFunc
	RespondWithJSON  http.HandlerFunc
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	type errResp struct {
		Error string `json:"error"`
	}
	if code > 499 {
		log.Printf("responding with 5xx error: %s", msg)
	}
	RespondWithJSON(w, code, errResp{
		Error:  msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Unable to marshal payload: %v", err)
		w.WriteHeader(500)
		return 
	}
	w.WriteHeader(code)
	w.Write(data)
	return 
}