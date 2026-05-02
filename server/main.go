package main

import (
	"log"
	"os"

	"leblanc/server/internal/db"
	"leblanc/server/internal/graph"
	"leblanc/server/internal/handlers"
	"leblanc/server/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	db.Init()
	services.EnsureAdminUser()

	r := gin.Default()

	// Configure CORS to allow requests from frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:3000", "http://127.0.0.1:5173", "https://le-blanc-web.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	// REST API endpoints
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "LeBlanc Go API with REST & GraphQL."}) })
	r.GET("/drinks", handlers.GetDrinks)
	r.POST("/reco/from-features", handlers.RecoFromFeatures)
	r.POST("/bookings", handlers.CreateBooking)
	r.POST("/auth/register", handlers.RegisterUser)
	r.POST("/auth/login", handlers.LoginUser)
	r.POST("/auth/request-verify", handlers.RequestVerify)
	r.POST("/auth/verify", handlers.VerifyToken)

	admin := r.Group("/admin")
	admin.Use(handlers.RequireAdminAuth())
	{
		admin.GET("/users", handlers.AdminListUsers)
		admin.POST("/users", handlers.AdminCreateUser)
		admin.PUT("/users/:id", handlers.AdminUpdateUser)
		admin.DELETE("/users/:id", handlers.AdminDeleteUser)

		admin.GET("/bookings", handlers.AdminListBookings)
		admin.POST("/bookings", handlers.AdminCreateBooking)
		admin.PUT("/bookings/:id", handlers.AdminUpdateBooking)
		admin.DELETE("/bookings/:id", handlers.AdminDeleteBooking)

		admin.GET("/drinks", handlers.AdminListDrinks)
		admin.POST("/drinks", handlers.AdminCreateDrink)
		admin.PUT("/drinks/:id", handlers.AdminUpdateDrink)
		admin.DELETE("/drinks/:id", handlers.AdminDeleteDrink)
	}

	// GraphQL endpoint
	r.POST("/graphql", graph.Handler())
	r.GET("/graphql", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "GraphQL endpoint - send POST requests with GraphQL queries"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	log.Println("Server listening on http://localhost:" + port)
	log.Println("REST API: http://localhost:" + port)
	log.Println("GraphQL: http://localhost:" + port + "/graphql")
	_ = r.Run(":" + port)
}
