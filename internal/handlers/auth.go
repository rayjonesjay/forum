package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"forum/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db DatabaseInterface
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *models.User `json:"user,omitempty"`
}

func NewAuthHandler(db DatabaseInterface) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendAuthResponse(w, false, "Invalid request format", nil)
		return
	}

	// Validate input
	if req.Email == "" || req.Username == "" || req.Password == "" {
		h.sendAuthResponse(w, false, "All fields are required", nil)
		return
	}

	// Check if email already exists
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", req.Email).Scan(&exists)
	if err != nil {
		h.sendAuthResponse(w, false, "Database error", nil)
		return
	}
	if exists {
		h.sendAuthResponse(w, false, "Email already registered", nil)
		return
	}

	// Check if username already exists
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", req.Username).Scan(&exists)
	if err != nil {
		h.sendAuthResponse(w, false, "Database error", nil)
		return
	}
	if exists {
		h.sendAuthResponse(w, false, "Username already taken", nil)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.sendAuthResponse(w, false, "Failed to process password", nil)
		return
	}

	// Insert user
	result, err := h.db.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)",
		req.Email, req.Username, string(hashedPassword))
	if err != nil {
		h.sendAuthResponse(w, false, "Failed to create user", nil)
		return
	}

	userID, _ := result.LastInsertId()
	user := &models.User{
		ID:       int(userID),
		Email:    req.Email,
		Username: req.Username,
		Created:  time.Now(),
	}

	h.sendAuthResponse(w, true, "User registered successfully", user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendAuthResponse(w, false, "Invalid request format", nil)
		return
	}

	// Get user from database
	var user models.User
	var hashedPassword string
	err := h.db.QueryRow("SELECT id, email, username, password, created FROM users WHERE email = ?", req.Email).
		Scan(&user.ID, &user.Email, &user.Username, &hashedPassword, &user.Created)
	if err != nil {
		h.sendAuthResponse(w, false, "Invalid credentials", nil)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		h.sendAuthResponse(w, false, "Invalid credentials", nil)
		return
	}

	// Create session
	sessionID, err := h.generateSessionID()
	if err != nil {
		h.sendAuthResponse(w, false, "Failed to create session", nil)
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour) // 24 hours
	_, err = h.db.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, user.ID, expiresAt)
	if err != nil {
		h.sendAuthResponse(w, false, "Failed to create session", nil)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expiresAt,
		HttpOnly: true,
		Path:     "/",
	})

	h.sendAuthResponse(w, true, "Login successful", &user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Delete session from database
		h.db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})

	h.sendAuthResponse(w, true, "Logged out successfully", nil)
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := h.getCurrentUser(r)
	if user == nil {
		h.sendAuthResponse(w, false, "Not authenticated", nil)
		return
	}

	h.sendAuthResponse(w, true, "User found", user)
}

func (h *AuthHandler) getCurrentUser(r *http.Request) *models.User {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil
	}

	var user models.User
	err = h.db.QueryRow(`
		SELECT u.id, u.email, u.username, u.created 
		FROM users u 
		JOIN sessions s ON u.id = s.user_id 
		WHERE s.id = ? AND s.expires_at > CURRENT_TIMESTAMP`,
		cookie.Value).Scan(&user.ID, &user.Email, &user.Username, &user.Created)
	
	if err != nil {
		return nil
	}

	return &user
}

func (h *AuthHandler) generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (h *AuthHandler) sendAuthResponse(w http.ResponseWriter, success bool, message string, user *models.User) {
	w.Header().Set("Content-Type", "application/json")
	
	if !success {
		w.WriteHeader(http.StatusBadRequest)
	}

	response := AuthResponse{
		Success: success,
		Message: message,
		User:    user,
	}

	json.NewEncoder(w).Encode(response)
}