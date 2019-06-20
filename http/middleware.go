package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type AuthRequest struct {
	SessionID string `json:"session_id"`
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try copy a request to read it twice.
		var buf bytes.Buffer
		tee := io.TeeReader(r.Body, &buf)
		newRequest, err := http.NewRequest(r.Method, r.URL.String(), tee)
		if err != nil {
			log.Printf("HTTP: Failed to copy request: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		decoder := json.NewDecoder(bytes.NewReader(buf.Bytes()))
		authRequest := new(AuthRequest)
		if err := decoder.Decode(authRequest); err != nil || authRequest.SessionID == "" {
			ctx := context.WithValue(r.Context(), "isLoggedIn", false)
			next.ServeHTTP(w, newRequest.WithContext(ctx))

			return
		}

		user, err := h.UserSession.Authenticate(authRequest.SessionID)
		if err != nil {
			log.Printf("HTTP: User not authenticated: %s", err)

			ctx := context.WithValue(r.Context(), "isLoggedIn", false)
			next.ServeHTTP(w, newRequest.WithContext(ctx))

			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "isLoggedIn", true)
		ctx = context.WithValue(ctx, "sessionID", authRequest.SessionID)
		next.ServeHTTP(w, newRequest.WithContext(ctx))
	})
}
