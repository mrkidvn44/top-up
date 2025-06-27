package util

import (
	"bytes"
	"net/http"
)

func SendPostRequest(url string, payload []byte) {
	http.Post(url, "application/json", bytes.NewBuffer(payload))
}
