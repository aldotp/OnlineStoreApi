package helper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Status  string `json:"status,omitempty"`
	Error   bool   `json:"error,omitempty"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewResponse() *Response {
	return &Response{}
}

// READ JSON
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// WRITE JSON
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJSON(response any, w http.ResponseWriter, status ...int) error {
	statusCode := http.StatusBadRequest // default

	if len(status) > 0 {
		statusCode = status[0]
	}

	return WriteJSON(w, statusCode, response)
}
