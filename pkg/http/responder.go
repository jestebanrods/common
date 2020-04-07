package http

import (
	"encoding/json"
	"net/http"
)

type HttpResponser interface {
	String(w http.ResponseWriter, code int, data string) error
	NoContent(w http.ResponseWriter, code int) error
	JSON(w http.ResponseWriter, code int, data interface{}) error
	JSONBlob(w http.ResponseWriter, code int, data []byte) error
}

type httpResponder struct{}

func (r *httpResponder) String(w http.ResponseWriter, code int, data string) error {
	w.WriteHeader(code)
	_, err := w.Write([]byte(data))
	return err
}

func (r *httpResponder) NoContent(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}

func (r *httpResponder) JSON(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func (r *httpResponder) JSONBlob(w http.ResponseWriter, code int, data []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(data)
	return err
}

func NewHttpResponder() *httpResponder {
	return &httpResponder{}
}
