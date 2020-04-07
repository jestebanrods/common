package http

import (
	"net/http"
)

type notFoundRedirectWrapper struct {
	http.ResponseWriter
	status int
}

func (w *notFoundRedirectWrapper) WriteHeader(status int) {
	w.status = status
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *notFoundRedirectWrapper) Write(p []byte) (int, error) {
	if w.status != http.StatusNotFound {
		return w.ResponseWriter.Write(p)
	}

	return len(p), nil
}

func notFoundWrapper(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nfrw := &notFoundRedirectWrapper{ResponseWriter: w}
		h.ServeHTTP(nfrw, r)

		if nfrw.status == 404 {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}
