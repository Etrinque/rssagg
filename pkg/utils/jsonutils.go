package utils

// Json respose for endpoint
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type JsonUtils struct {
	RespondWithError http.HandlerFunc
	RespondWithJSON  http.HandlerFunc
}

// helper func for json logic parsing
func JSONdecode[T any](r *http.Request) (T, error) {
    var v T
    err := json.NewDecoder(r.Body).Decode(v); if err != nil {
		return v, err
    }
    return v, nil
}

func JSONencode[T any](w http.ResponseWriter, r *http.Request,status int,v T) error {
    w.Header().Set("Content-Type", "aplication/json")
    w.WriteHeader(status)
    err := json.NewEncoder(w).Encode(v); if err != nil {
        errors.New("Err: could not encode json")
		return err
    }
    return nil
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