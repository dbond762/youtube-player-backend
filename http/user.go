package http

import (
	"encoding/json"
	"log"
	"net/http"

	player "github.com/dbond762/youtube-player-backend"
)

type LoginRequest struct {
	User struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	} `json:"user"`
}

type LoginResponse struct {
	User struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		SessionID string `json:"session_id"`
		Likes     []struct {
			VideoID int64 `json:"video_id"`
		} `json:"likes"`
	} `json:"user"`
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var loginRequest LoginRequest
	if err := decoder.Decode(&loginRequest); err != nil {
		log.Printf("Error on decode login request: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user := &player.User{
		Login:    loginRequest.User.Login,
		Password: loginRequest.User.Password,
	}

	user, err := h.UserService.CreateUser(*user)
	if err != nil {
		log.Printf("Error on create user: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, sessionID, err := h.UserSession.Login(user.Login)
	if err != nil {
		log.Printf("Error on login user: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var loginResp LoginResponse
	loginResp.User.ID = user.ID
	loginResp.User.Login = user.Login
	loginResp.User.SessionID = sessionID

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&loginResp); err != nil {
		log.Printf("Error on encoding login result: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	loginForm := new(LoginRequest)
	if err := decoder.Decode(loginForm); err != nil {
		log.Printf("Error on decode login request: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, sessionID, err := h.UserSession.Login(loginForm.User.Login)
	if err != nil {
		log.Printf("Error on login user: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var loginResult LoginResponse
	loginResult.User.ID = user.ID
	loginResult.User.Login = user.Login
	loginResult.User.SessionID = sessionID

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&loginResult); err != nil {
		log.Printf("Error on encoding login result: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID, ok := ctx.Value("sessionID").(string)
	if !ok {
		log.Printf("Not found session_id cookie")
		return
	}

	h.UserSession.Logout(sessionID)
}
