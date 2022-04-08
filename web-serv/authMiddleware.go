package main

import (
	"context"
	"net/http"
)

func authRequiredMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO actually do authentication
		w.Header().Add("AuchCheckerVersion", "FAKE")
		ctx := context.WithValue(r.Context(), "email", "fake@fake_email.fake")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
