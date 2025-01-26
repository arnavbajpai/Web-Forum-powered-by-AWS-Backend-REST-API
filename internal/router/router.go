package router

import (
	"github.com/arnavbajpai/web-forum-project/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
)
func injectTestContextMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("userEmail", "arnav@example.com")
        c.Set("userGivenName", "arnav")
        c.Set("userFamilyName", "bajpai")
        c.Next()
    }
}
func Setup() chi.Router {
	r := chi.NewRouter()
	setUpRoutes(r)
	return r
}

func setUpRoutes(r chi.Router) {
	r.Group(routes.GetRoutes())
}
