package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/middleware"
	"backend/internal/service"
)

type Handler struct {
	userSvc *service.UserService
	txSvc   *service.TransactionService
	balSvc  *service.BalanceService
}

func NewHandler(u *service.UserService, t *service.TransactionService, b *service.BalanceService) *Handler {
	return &Handler{userSvc: u, txSvc: t, balSvc: b}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.userSvc.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, token, err := h.userSvc.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromUserID *int64 `json:"from_user_id"`
		ToUserID   *int64 `json:"to_user_id"`
		Amount     int64  `json:"amount"`
		Type       string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tx, err := h.txSvc.Create(r.Context(), req.FromUserID, req.ToUserID, req.Amount, req.Type)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusAccepted, tx)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
    respondError(w, http.StatusNotImplemented, "Refresh not implemented yet")
}

func (h *Handler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
    userIDVal := r.Context().Value(middleware.UserIDKey)
    if userIDVal == nil {
         respondError(w, http.StatusUnauthorized, "Unauthorized")
         return
    }
    userID := int64(userIDVal.(float64))
    
    txs, err := h.txSvc.GetHistory(r.Context(), userID)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondJSON(w, http.StatusOK, txs)
}

func (h *Handler) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {
    userIDVal := r.Context().Value(middleware.UserIDKey)
    if userIDVal == nil {
         respondError(w, http.StatusUnauthorized, "Unauthorized")
         return
    }
    userID := int64(userIDVal.(float64))
    
    logs, err := h.balSvc.GetHistory(r.Context(), userID)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondJSON(w, http.StatusOK, logs)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.userSvc.ListUsers(r.Context())
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondJSON(w, http.StatusOK, users)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        respondError(w, http.StatusBadRequest, "Invalid ID")
        return
    }
    err = h.userSvc.DeleteUser(r.Context(), id)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value(middleware.UserIDKey)
    if userIDVal == nil {
         respondError(w, http.StatusUnauthorized, "Unauthorized")
         return
    }
    userID := int64(userIDVal.(float64))

	bal, err := h.balSvc.GetBalance(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch balance")
		return
	}
	respondJSON(w, http.StatusOK, bal)
}

