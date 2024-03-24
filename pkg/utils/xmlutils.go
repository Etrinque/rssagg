package utils

import (
	"encoding/xml"
	"fmt"
	"net/http"
)


func XMLdecode[T any](r *http.Request) (T, error) {
	var v T
	err := xml.NewDecoder(r.Body).Decode(v); if err != nil {
		 return v, err
	}
	return v, nil
 }
 
func XMLencode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	 w.Header().Set("Content-Type", "application/xml")
	 w.WriteHeader(status)
	 
	 err := xml.NewEncoder(w).Encode(v); if err != nil {
		 fmt.Println("Err: could not encode xml")
	 }
	 return nil
 }