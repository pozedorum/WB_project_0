// Package middleware is responcible for
package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"wb_project_0/credentials"
)

func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			askForCredentials(w)
			return
		}

		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			http.Error(w, "Invalid auth header", http.StatusBadRequest)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			http.Error(w, "Invalid base64 encoding", http.StatusBadRequest)
			return
		}

		newCreds := strings.SplitN(string(decoded), ":", 2)
		if len(newCreds) != 2 {
			askForCredentials(w)
			return
		}

		if newCreds[0] != credentials.Creds.Login || newCreds[1] != credentials.Creds.Password {
			askForCredentials(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func askForCredentials(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Admin Access"`)
	http.Error(w, "Authentication required", http.StatusUnauthorized)
}
