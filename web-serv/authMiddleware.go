package main

import "net/http"

func authRequiredMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO actually do authentication
		next.ServeHTTP(w, r)
	})
}
