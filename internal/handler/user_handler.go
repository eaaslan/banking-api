package handler

import (
	"encoding/json"
	"go-banking-api/internal/domain"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Role     *string `json:"role,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	user, err := h.userService.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusCreated, user, "User registered successfully")
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	user, err := h.userService.Login(r.Context(), req.UsernameOrEmail, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, user, "Login successful")
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), uint(id))
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, user, "")
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	users, err := h.userService.GetAllUsers(r.Context(), page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, users, "")
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	updates := make(map[string]interface{})
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}

	user, err := h.userService.UpdateUser(r.Context(), uint(id), updates)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, user, "User updated successfully")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.userService.DeleteUser(r.Context(), uint(id))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, nil, "User deleted successfully")
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	idStr = strings.TrimSuffix(idStr, "/change-password")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	err = h.userService.ChangePassword(r.Context(), uint(id), req.OldPassword, req.NewPassword)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, nil, "Password changed successfully")
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
	slog.Warn("Request failed", "status", statusCode, "message", message)
}

func respondSuccess(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}
