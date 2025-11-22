package api

import (
	"md6/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Gin Demo API",
			"endpoints": gin.H{
				"register":        "/register (POST)",
				"login":           "/login (POST)",
				"change-password": "/change-password (POST)",
				"users":           "/users (GET) - Debug only",
				"verify-token":    "/verify-token (POST)",
				"profile":         "/profile (GET) - Protected",
			},
		})
	})
	r.POST("/register", register)
	r.POST("/login", login)
	r.POST("/change-password", changePassword)
	r.GET("/users", getUsers)
	r.POST("/verify-token", verifyToken)
	protected := r.Group("/")
	protected.Use(utils.JWTAuthMiddleware())
	{
		protected.GET("/profile", getUserProfile)
	}
	r.GET("/login-page", func(c *gin.Context) {
		c.File("./static/login.html")
	})
	r.Run(":8088")
}
