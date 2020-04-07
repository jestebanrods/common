package http

import (
	"encoding/json"
	"net/http"
)

type Responder interface {
	String(code int, data string) error
	NoContent(code int) error
	JSON(code int, data interface{}) error
	JSONBlob(code int, data []byte) error
}

type responder struct {
	writer http.ResponseWriter
}

func (r *responder) String(code int, data string) error {
	r.writer.WriteHeader(code)
	_, err := r.writer.Write([]byte(data))
	return err
}

func (r *responder) NoContent(code int) error {
	r.writer.WriteHeader(code)
	return nil
}

func (r *responder) JSON(code int, data interface{}) error {
	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.WriteHeader(code)
	return json.NewEncoder(r.writer).Encode(data)
}

func (r *responder) JSONBlob(code int, data []byte) error {
	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.WriteHeader(code)
	_, err := r.writer.Write(data)
	return err
}

func NewResponder(w http.ResponseWriter) *responder {
	return &responder{writer: w}
}
