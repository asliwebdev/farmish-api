package handlers

import (
	"farmish/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary		Create Farm
// @Description	Create Farm for owner.
// @Tags			farms
// @Accept application/json
// @Produce application/json
// @Param request body models.CreateFarmRequest true "Request body for creating a farm"
// @Success 201 {object} models.CreateFarmResponse
// @Failure		400		{object}	models.ErrResp	"Invalid request body"
// @Failure		500		{object}	models.ErrResp	"Internal server error"
// @Security		BearerAuth
// @Router			/farms [post]
func (h *Handler) CreateFarm(c *gin.Context) {
	var farm models.Farm
	if err := c.ShouldBindJSON(&farm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.farmService.CreateFarm(&farm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Farm created successfully", "farm": farm})
}

// @Summary		Get a farm by ID
// @Description	Retrieve farm details by their UUID.
// @Tags			farms
// @Produce		application/json
// @Param			id		path		string			true	"Farm ID (UUID)"
// @Success		200		{object}	models.Farm		"Farm details"
// @Failure		400		{object}	models.ErrResp	"Invalid farm ID format"
// @Failure		404		{object}	models.ErrResp	"Farm not found"
// @Failure		500		{object}	models.ErrResp	"Internal server error"
// @Security		BearerAuth
// @Router			/farms/{id} [get]
func (h *Handler) GetFarmByID(c *gin.Context) {
	id := c.Param("id")
	farmID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Farm ID is not in the UUID format"})
		return
	}

	farm, err := h.farmService.GetFarmByID(farmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if farm == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "farm with this ID not found"})
		return
	}

	c.JSON(http.StatusOK, farm)
}

// @Summary		Get all farms
// @Description	Retrieve a list of all existing farms.
// @Tags			farms
// @Produce		application/json
// @Success		200		{array}		models.Farm		"List of farms"
// @Failure		500		{object}	models.ErrResp	"Internal server error"
// @Security		BearerAuth
// @Router			/farms [get]
func (h *Handler) GetAllFarms(c *gin.Context) {
	farms, err := h.farmService.GetAllFarms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, farms)
}

// @Summary		Update a farm
// @Description	Update details of a farm by their UUID.
// @Tags			farms
// @Accept			application/json
// @Produce		application/json
// @Param			id		path		string			true	"Farm ID (UUID)"
// @Param			request	body		models.UpdateFarmRequest	true	"Farm update payload"
// @Success		200		{object}	models.UpdateFarmResp			"Farm updated successfully"
// @Failure		400		{object}	models.ErrResp	"Invalid input or ID format"
// @Failure		500		{object}	models.ErrResp	"Internal server error"
// @Security		BearerAuth
// @Router			/farms/{id} [put]
func (h *Handler) UpdateFarm(c *gin.Context) {
	var farm models.UpdateFarmRequest
	if err := c.ShouldBindJSON(&farm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.farmService.UpdateFarm(&farm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Farm updated successfully", "farm": farm})
}

// @Summary		Delete a farm
// @Description	Delete a farm by their UUID.
// @Tags			farms
// @Produce		application/json
// @Param			id		path		string			true	"Farm ID (UUID)"
// @Success		200		{object}	models.MessageResp			"Farm deleted successfully"
// @Failure		400		{object}	models.ErrResp	"Invalid Farm ID format"
// @Failure		500		{object}	models.ErrResp	"Internal server error"
// @Security		BearerAuth
// @Router			/farms/{id} [delete]
func (h *Handler) DeleteFarm(c *gin.Context) {
	id := c.Param("id")
	farmID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Farm ID is not in the UUID format"})
		return
	}

	if err := h.farmService.DeleteFarm(farmID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Farm deleted successfully"})
}
