package main

import (
	"fmt"

	"github.com/arnavbajpai/web-forum-project/internal/api"
	"github.com/arnavbajpai/web-forum-project/internal/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func injectTestContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userEmail", "arnav@example.com")
		c.Set("userGivenName", "arnav")
		c.Set("userFamilyName", "bajpai")
		c.Next()
	}
}
func main() {
	err := database.InitializeDB()
	if err != nil {
		fmt.Println("Connection error.")
	}
	fmt.Println("Connected!")
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	initializeRoutes(r)
	r.Run(":8080")
}

func initializeRoutes(r *gin.Engine) {
	r.Use(injectTestContextMiddleware())
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", api.AddUser)
		userRoutes.GET("", api.GetUserByAlias)
		userRoutes.POST("/:userID", api.UpdateUser)
		userRoutes.DELETE("/:userID", api.DeleteUser)
		userRoutes.GET("/:userID", api.GetUserByID)
	}

	postRoutes := r.Group("/posts")
	{
		postRoutes.POST("", api.AddPost)
		postRoutes.GET("", api.GetPosts)
		postRoutes.POST("/:postID", api.UpdatePost)
		postRoutes.DELETE("/:postID", api.DeletePost)
		postRoutes.GET("/:postID", api.GetPostByID)
	}

	categoryRoutes := r.Group("/categories")
	{
		categoryRoutes.GET("", api.GetCategories)
		categoryRoutes.POST("", api.AddCategory)
	}

	tagRoutes := r.Group("/tags")
	{
		tagRoutes.GET("", api.GetTopTags)
		tagRoutes.DELETE("/:tagID", api.DeleteTag) // CRON
		tagRoutes.GET("/:tagID", api.GetTag)
	}
}
