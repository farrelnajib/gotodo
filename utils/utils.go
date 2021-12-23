package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var Message = func(status string, message string, data interface{}) Response {
	return Response{Status: status, Message: message, Data: data}
}

var Respond = func(w http.ResponseWriter, code int, response Response) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := response
	json.NewEncoder(w).Encode(resp)
}
