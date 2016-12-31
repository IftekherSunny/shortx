package http

import (
	"encoding/json"
	netHttp "net/http"
)

// Response struct
type ResponseStruct struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"status_code"`
}

// Reponse method
func Response(w netHttp.ResponseWriter, statuscode int, data interface{}) {
	rs := ResponseStruct{
		Data:       data,
		StatusCode: statuscode,
	}

	js, _ := json.Marshal(rs)

	w.Header().Set("content-type", "application/json")
	w.Write(js)
}
