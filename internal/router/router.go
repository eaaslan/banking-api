package router

import (
	"go-banking-api/internal/handler"
	"net/http"
)

func SetupRoutes(userHandler *handler.UserHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/auth/register", userHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", userHandler.Login)

	mux.HandleFunc("/api/v1/users", userHandler.GetAllUsers)
	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/users/" {
			userHandler.GetAllUsers(w, r)
			return
		}

		//if r.URL.Path[len(r.URL.Path)-16:] == "/change-password" {
		//	userHandler.ChangePassword(w, r)
		//	return
		//}

		switch r.Method {
		case http.MethodGet:
			userHandler.GetUser(w, r)
		case http.MethodPut, http.MethodPatch:
			userHandler.UpdateUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}
