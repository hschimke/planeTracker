package main

import (
	"context"
	"net/http"
	"strings"
)

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-PlaneTracker-Auth-Type-Request")
}

func authRequiredMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}

		// TODO actually do authentication
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			http.Error(w, "Malformed Token", http.StatusUnauthorized)
			return
		}
		authTypeRequest := r.Header.Get("X-PlaneTracker-Auth-Type-Request")

		w.Header().Add("AuchCheckerVersion", authTypeRequest)
		ctx := context.WithValue(r.Context(), "email", "fake@fake.fake")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
