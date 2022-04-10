package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
)

func authGoogleToken(ctx context.Context, tokenString string, audience string) (string, error) {
	token, validateErr := idtoken.Validate(ctx, tokenString, audience)
	if validateErr != nil {
		return "", validateErr
	}
	if token.Claims["email_verified"].(string) != "true" {
		return "", fmt.Errorf("unverified email")
	}
	return token.Claims["email"].(string), nil
}

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

		var ctx context.Context

		switch authTypeRequest {
		case "fake":
			ctx = context.WithValue(r.Context(), "email", "fake@fake.fake")
		case "google":
			// TODO add audience
			email, authErr := authGoogleToken(r.Context(), authHeader[1], "")
			if authErr != nil {
				http.Error(w, "google auth failed, token appears fake", http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(r.Context(), "email", email)
		case "microsoft":
			// TODO validate microsoft token
			http.Error(w, "microsoft jwt not supported", http.StatusUnauthorized)
			return
		default:
			http.Error(w, fmt.Sprintf("authtype: '%s' jwt not supported", authTypeRequest), http.StatusUnauthorized)
			return
		}

		w.Header().Add("AuchCheckerVersion", authTypeRequest)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
