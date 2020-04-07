package http

import (
	"net/http"
)

type NotFoundRedirectWrapper struct {
	http.ResponseWriter
	status int
}

func (w *NotFoundRedirectWrapper) WriteHeader(status int) {
	w.status = status
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *NotFoundRedirectWrapper) Write(p []byte) (int, error) {
	if w.status != http.StatusNotFound {
		return w.ResponseWriter.Write(p)
	}

	return len(p), nil
}

func NotFoundWrapper(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrapper := &NotFoundRedirectWrapper{ResponseWriter: w}
		h.ServeHTTP(wrapper, r)

		if wrapper.status == 404 {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}
