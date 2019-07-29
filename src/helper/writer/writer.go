// Package writer : write http response in JSON
package writer

import (
	"encoding/json"
	"log"
	"net/http"
)

// Writer wrapping http.ResponseWriter
type Writer struct {
	Writer http.ResponseWriter
}

// New : Create new writer from http.ResponseWriter
func New(w http.ResponseWriter) *Writer {
	return &Writer{
		Writer: w,
	}
}

// Success response
func (w *Writer) Success(data interface{}) {
	// Success response format
	resp := struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		"000",
		"Success",
		data,
	}

	w.writeJSON(resp, http.StatusOK)
}

// Error response
func (w *Writer) Error(err error) {
	// Error response format
	data := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		"001",
		err.Error(),
	}

	w.writeJSON(data, http.StatusInternalServerError)
}

// Write JSON response
func (w *Writer) writeJSON(data interface{}, status int) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	w.Writer.Header().Set("Content-Type", "application/json")
	w.Writer.WriteHeader(status)
	w.Writer.Write(b)
}
