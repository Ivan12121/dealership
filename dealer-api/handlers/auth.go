package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("super-secret-key")

type AuthHandler struct {
	DB *sql.DB
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if req.Password != req.ConfirmPassword {
		http.Error(w, "Password not equals", http.StatusBadRequest)
		return
	}

	hashedPAssword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error generate password hash", http.StatusInternalServerError)
		return
	}

	var userID int
	err = h.DB.QueryRow(
		"INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, 'client') RETURNING id",
		req.Name, req.Email, string(hashedPAssword),
	).Scan(&userID)

	if err != nil {
		http.Error(w, "User already exist", http.StatusConflict)
		return
	}

	token, err := generateJwt(userID, req.Email, "client")
	if err != nil {
		http.Error(w, "Error generate JWT token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token":token})
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var userID int
	var storeHash, role string

	err := h.DB.QueryRow(
		"SELECT id, password_hash, role FROM users WHERE email = $1", 
		req.Email,
	).Scan(&userID, &storeHash, &role)

	if err == sql.ErrNoRows {
		http.Error(w, "Wrong email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "DataBase Error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storeHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Wrong email or password", http.StatusUnauthorized)
		return
	}

	token, err := generateJwt(userID, req.Email, role)
	if err != nil {
		http.Error(w, "Error generate JWT token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token":token})
}

func generateJwt(userID int, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email": email,
		"role": role,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}