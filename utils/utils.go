package utils

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var Message = func(status string, message string, data interface{}) Response {
	return Response{Status: status, Message: message, Data: data}
}

var Respond = func(c *fasthttp.RequestCtx, code int, response Response) {
	c.Response.Header.SetContentType("application/json")
	c.Response.SetStatusCode(code)
	json.NewEncoder(c).Encode(response)
}
