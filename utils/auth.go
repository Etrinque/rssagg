package utils

import (
	"errors"
	"net/http"
	"strings"
)

var ErrNotAuthorized = errors.New("no authorization header")

func GetApiToken(headers http.Header) (string, error) {

	header := headers.Get("Authorization")
	if header == "" {
		return "", ErrNotAuthorized
	}
	splitHeader := strings.Split(header, " ")
	if len(splitHeader) < 2 || splitHeader[0] != "ApiKey" {
		return "", errors.New("malformed header")
	}
	return splitHeader[1], nil
}
