package handlers

import (
	"farmish/internal/services"
	"farmish/pkg/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	userService *services.UserService
}

func NewHandler(userService *services.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

// Run ...
// @title           Farmish API
// @version         1.0
// @description     Testing Swagger APIs.
// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io
// @host            localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in				header
// @name			Authorization
// @type 			apikey
// @schema 			bearer
// @bearerFormat	JWT
func Run(h *Handler) *gin.Engine {
	router := gin.Default()

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// AUTH ROUTES
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", h.Login)
		authRoutes.POST("/signup", h.SignUp)
	}

	router.Use(middleware.AuthMiddleware())

	// USER ROUTES
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", h.GetAllUsers)
		userRoutes.GET("/:id", h.GetUserByID)
		userRoutes.PUT("/:id", h.UpdateUser)
		userRoutes.DELETE("/:id", h.DeleteUser)
	}

	return router
}
