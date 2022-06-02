package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response is a struct response
type Response struct {
	Data interface{} `json:"data"`
}

// Empty is a empty struct response
type Empty struct{}

// ResponseJSON return a response with status code
func ResponseJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data == nil {
		data = Empty{}
	}
	var response = Response{
		Data: data,
	}
	responseWriter(w, &response)
}

// responseWriter return a response
func responseWriter(w http.ResponseWriter, response *Response) {
	responseData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something went wrong :(")
	}
	fmt.Fprint(w, string(responseData))
}
