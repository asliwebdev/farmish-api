package handlers

import (
	"farmish/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create a new animal
// @Description Add a new animal to a farm
// @Tags animals
// @Accept application/json
// @Produce application/json
// @Param request body models.CreateAnimalReq true "Animal data"
// @Success 201 {object} models.CreateAnimalResp
// @Failure 400 {object} models.ErrResp "Invalid input"
// @Failure 500 {object} models.ErrResp "Internal server error"
// @Security BearerAuth
// @Router /animals [post]
func (h *Handler) CreateAnimal(c *gin.Context) {
	var animal models.AnimalWithoutTime
	if err := c.ShouldBindJSON(&animal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.animalService.CreateAnimal(&animal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if animal.HealthStatus == "" {
		animal.HealthStatus = "Healthy"
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Animal added to farm successfully", "animal": animal})
}

// @Summary Get an animal by ID
// @Description Retrieve details of an animal using its unique ID
// @Tags animals
// @Produce application/json
// @Param id path string true "Animal ID(UUID)"
// @Success 200 {object} models.Animal
// @Failure 400 {object} models.ErrResp "Invalid animal ID"
// @Failure 404 {object} models.ErrResp "Animal not found"
// @Failure 500 {object} models.ErrResp "Internal server error"
// @Security BearerAuth
// @Router /animals/{id} [get]
func (h *Handler) GetAnimalByID(c *gin.Context) {
	idParam := c.Param("id")
	animalID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	animal, err := h.animalService.GetAnimalByID(animalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if animal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "animal with this ID not found"})
		return
	}

	c.JSON(http.StatusOK, animal)
}

// @Summary Get animals by farm ID
// @Description Retrieve a list of animals for a specific farm
// @Tags animals
// @Produce application/json
// @Param farm_id query string true "Farm ID"
// @Success 200 {array} models.Animal
// @Failure 400 {object} models.ErrResp "Invalid farm ID"
// @Failure 500 {object}  models.ErrResp "Internal server error"
// @Security BearerAuth
// @Router /animals [get]
func (h *Handler) GetAnimalsByFarmID(c *gin.Context) {
	farmIDParam := c.Query("farm_id")
	farmID, err := uuid.Parse(farmIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid farm ID"})
		return
	}

	animals, err := h.animalService.GetAnimalsByFarmID(farmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animals)
}

// @Summary Update an animal
// @Description Modify details of an existing animal
// @Tags animals
// @Accept json
// @Produce json
// @Param request body models.UpdateAnimalReq true "Updated animal data"
// @Success 200 {object} models.UpdateAnimalResp
// @Failure 400 {object} models.ErrResp "Invalid input"
// @Failure 500 {object} models.ErrResp "Internal server error"
// @Security BearerAuth
// @Router /animals [put]
func (h *Handler) UpdateAnimal(c *gin.Context) {
	var animal models.UpdateAnimalReq
	if err := c.ShouldBindJSON(&animal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.animalService.UpdateAnimal(&animal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Animal updated successfully", "animal": animal})
}

// @Summary Delete an animal
// @Description Remove an animal from the system
// @Tags animals
// @Produce application/json
// @Param id path string true "Animal ID"
// @Success		200		{object}	models.MessageResp	"Animal deleted successfully"
// @Failure 500 {object} models.ErrResp "Internal server error"
// @Security BearerAuth
// @Router /animals/{id} [delete]
func (h *Handler) DeleteAnimal(c *gin.Context) {
	idParam := c.Param("id")
	animalID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	if err := h.animalService.DeleteAnimal(animalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Animal deleted successfully"})
}
