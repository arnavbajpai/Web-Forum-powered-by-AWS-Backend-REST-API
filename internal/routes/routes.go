package routes

import (
	"encoding/json"
	"net/http"

	"github.com/arnavbajpai/web-forum-project/internal/handlers/users"
	"github.com/go-chi/chi/v5"
)

func GetRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/users", func(w http.ResponseWriter, req *http.Request) {
			response, _ := users.HandleUserList(w, req)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
	}
}
