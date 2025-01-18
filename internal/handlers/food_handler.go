package handlers

import (
	"farmish/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create a food item and add it to the warehouse
// @Description Create a new food item and associate it with a specific farm
// @Tags foods
// @Accept application/json
// @Produce application/json
// @Param request body models.AddFoodReq true "Warehouse Food"
// @Success 201 {object} models.AddFoodResp
// @Failure 400 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /foods [post]
func (h *Handler) AddFoodToWarehouse(c *gin.Context) {
	var food models.FoodWithoutTime
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.foodService.AddFoodToWarehouse(&food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Food added to the farm successfully", "food": food})
}

// @Summary Get all food items for a specific farm
// @Description Retrieve all foods associated with a specific farm
// @Tags foods
// @Produce application/json
// @Param farm_id path string true "Farm ID"
// @Success 200 {array} models.Food
// @Failure 400 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /foods/{farm_id} [get]
func (h *Handler) GetWarehouseFoods(c *gin.Context) {
	farmID, err := uuid.Parse(c.Param("farm_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid farm ID"})
		return
	}

	foods, err := h.foodService.GetFoodsByFarm(farmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, foods)
}

// @Summary Get a food by ID
// @Description Retrieve details of a food using its unique ID
// @Tags foods
// @Produce application/json
// @Param food_id path string true "Food ID(UUID)"
// @Success 200 {object} models.Food
// @Failure 400 {object} models.ErrResp "Invalid food ID"
// @Failure 404 {object} models.ErrResp "Food not found"
// @Failure 500 {object} models.ErrResp "Internal server error"
// @Security BearerAuth
// @Router /foods/food/{food_id} [get]
func (h *Handler) GetFoodByID(c *gin.Context) {
	foodID, err := uuid.Parse(c.Param("food_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID"})
		return
	}

	food, err := h.foodService.GetFoodByID(foodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if food == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "food not found with ID"})
		return
	}

	c.JSON(http.StatusOK, &food)
}

// @Summary Update a food item in the warehouse
// @Description Update an existing food item in the warehouse
// @Tags foods
// @Accept application/json
// @Produce application/json
// @Param request body models.UpdateFoodReq true "Warehouse Food"
// @Success 200 {object} models.UpdateFoodResp
// @Failure 400 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /foods [put]
func (h *Handler) UpdateWarehouseFood(c *gin.Context) {
	var food models.UpdateFoodReq
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.foodService.UpdateFood(&food); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "food updated successfully", "food": food})
}

// @Summary Remove a food item from the warehouse
// @Description Delete a food item from the warehouse
// @Tags foods
// @Produce application/json
// @Param food_id path string true "Food ID"
// @Success 200 {object} models.MessageResp
// @Failure 400 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /foods/{food_id} [delete]
func (h *Handler) RemoveWarehouseFood(c *gin.Context) {
	foodID, err := uuid.Parse(c.Param("food_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID"})
		return
	}

	if err := h.foodService.RemoveWarehouseFood(foodID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Food item removed successfully"})
}
