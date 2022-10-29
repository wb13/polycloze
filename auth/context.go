// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package auth

import (
	"context"
	"database/sql"
	"net/http"
)

type contextValueKey int

// Keys for getting values from request context.
const (
	keyUserDB contextValueKey = iota
	keySession
)

// Stuffs pointer to database of users into request context.
func Middleware(db *sql.DB) func(http.Handler) http.Handler {
	// Gets user session and stuffs it in the request context.
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), keyUserDB, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Gets pointer to user database from request context.
// Assumes Middleware is used.
func GetDB(r *http.Request) *sql.DB {
	return r.Context().Value(keyUserDB).(*sql.DB)
}
