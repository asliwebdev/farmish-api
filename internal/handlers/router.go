package handlers

import (
	_ "farmish/docs"
	"farmish/internal/services"
	"farmish/pkg/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	userService          *services.UserService
	farmService          *services.FarmService
	animalService        *services.AnimalService
	foodService          *services.FoodService
	medicineService      *services.MedicineService
	feedingRecordService *services.FeedingRecordService
	medicalRecordService *services.MedicalRecordService
}

func NewHandler(userService *services.UserService, farmService *services.FarmService,
	animalService *services.AnimalService,
	foodService *services.FoodService, medicineService *services.MedicineService,
	feedingRecordService *services.FeedingRecordService,
	medicalRecordService *services.MedicalRecordService,
) *Handler {
	return &Handler{
		userService:          userService,
		farmService:          farmService,
		animalService:        animalService,
		foodService:          foodService,
		medicineService:      medicineService,
		feedingRecordService: feedingRecordService,
		medicalRecordService: medicalRecordService,
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

	// FARM ROUTES
	farmRoutes := router.Group("/farms")
	{
		farmRoutes.POST("/", h.CreateFarm)
		farmRoutes.GET("/:id", h.GetFarmByID)
		farmRoutes.GET("/", h.GetAllFarms)
		farmRoutes.PUT("/:id", h.UpdateFarm)
		farmRoutes.DELETE("/:id", h.DeleteFarm)
	}

	// ANIMAL ROUTES
	animalRoutes := router.Group("/animals")
	{
		animalRoutes.POST("/", h.CreateAnimal)
		animalRoutes.GET("/:id", h.GetAnimalByID)
		animalRoutes.GET("/", h.GetAnimalsByFarmID)
		animalRoutes.PUT("/", h.UpdateAnimal)
		animalRoutes.DELETE("/:id", h.DeleteAnimal)
	}

	// FOOD ROUTES
	foodRoutes := router.Group("/foods")
	{
		foodRoutes.POST("/", h.AddFoodToWarehouse)
		foodRoutes.GET("/:farm_id", h.GetWarehouseFoods)
		foodRoutes.GET("/food/:food_id", h.GetFoodByID)
		foodRoutes.PUT("/", h.UpdateWarehouseFood)
		foodRoutes.DELETE("/:food_id", h.RemoveWarehouseFood)
	}

	// MEDICINE ROUTES
	medicineRoutes := router.Group("/medicines")
	{
		medicineRoutes.POST("/", h.CreateMedicine)
		medicineRoutes.GET("/", h.GetAllMedicines)
		medicineRoutes.GET("/:id", h.GetMedicineByID)
		medicineRoutes.PUT("/:id", h.UpdateMedicine)
		medicineRoutes.DELETE("/:id", h.DeleteMedicine)
	}

	// FEEDING RECORD ROUTES
	feedingRecordRoutes := router.Group("/feeding_records")
	{
		feedingRecordRoutes.POST("/", h.CreateFeedingRecord)
		feedingRecordRoutes.GET("/:id", h.GetFeedingRecordByID)
		feedingRecordRoutes.GET("/animal/:animal_id", h.GetFeedingRecordsByAnimalID)
		feedingRecordRoutes.PUT("/:id", h.UpdateFeedingRecord)
		feedingRecordRoutes.DELETE("/:id", h.DeleteFeedingRecord)
	}

	// MEDICAL RECORD ROUTES
	medicalRecords := router.Group("/medical_records")
	{
		medicalRecords.POST("", h.CreateMedicalRecord)
		medicalRecords.GET("/:id", h.GetMedicalRecordByID)
		medicalRecords.GET("/animals/:animal_id", h.GetMedicalRecordsByAnimalID)
		medicalRecords.PUT("/:id", h.UpdateMedicalRecord)
		medicalRecords.DELETE("/:id", h.DeleteMedicalRecord)
	}

	return router
}
