package routes

import (
	"good-guys/backend/api"
	"good-guys/backend/database"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func SetupAPIRoutes(router *gin.Engine) {
	setupCORS(router)

	handler := &api.Handler{DB: database.DB}
	store := memory.NewStore()

	loginLimiter := limiter.New(store, limiter.Rate{
		Period: 15 * time.Minute,
		Limit:  5,
	})

	authLimiter := limiter.New(store, limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  30,
	})

	apiRouter := router.Group("/api")
	{
		apiRouter.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello from Go!"})
		})
		apiRouter.GET("/", handler.GetEvents)
		apiRouter.POST("/event", handler.GetEvent)
		apiRouter.POST("/contribute", handler.ContributeEvent)

		// admin routes
		apiRouter.GET("/auth-me", mgin.NewMiddleware(authLimiter), handler.AuthMe)
		apiRouter.POST("/login", mgin.NewMiddleware(loginLimiter), handler.Login)
		adminRoutes := apiRouter.Group("/")
		adminRoutes.Use(handler.AdminMiddleware())
		{
			adminRoutes.GET("/admin-events", handler.GetAdminEvents)
		}
	}
}

func setupCORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow local frontend in dev
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func SetupStaticRoutes(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Next() // Let API routes handle it
			return
		}

		// Try to serve static files
		filePath := "./frontend/build" + c.Request.URL.Path
		if _, err := os.Stat(filePath); err == nil {
			c.File(filePath)
			c.Abort()
			return
		}
		// If not a file, serve index.html (for SvelteKit routing)
		c.File("./frontend/build/index.html")
		c.Abort()
	})
}
