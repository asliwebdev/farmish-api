package handlers

import (
	"farmish/internal/models"
	"farmish/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create a new medicine
// @Description Add a new medicine to the farm
// @Tags medicines
// @Accept application/json
// @Produce application/json
// @Param request body models.MedicineReq true "Medicine Details"
// @Success 201 {object} models.MedicineResp
// @Failure 400 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /medicines [post]
func (h *Handler) CreateMedicine(c *gin.Context) {
	var medicine models.MedicineWithoutTime
	if err := c.ShouldBindJSON(&medicine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.medicineService.CreateMedicine(&medicine); err != nil {
		if err == services.ErrQuantityLessThanThreshold {
			c.JSON(http.StatusBadRequest, gin.H{"error": services.ErrQuantityLessThanThreshold})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Medicine created successfully", "medicine": medicine})
}

// @Summary Get all medicines
// @Description Retrieve all medicines for a specific farm
// @Tags medicines
// @Param farm_id query string true "Farm ID"
// @Success 200 {array} models.Medicine
// @Failure 400 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /medicines [get]
func (h *Handler) GetAllMedicines(c *gin.Context) {
	farmID, err := uuid.Parse(c.Query("farm_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid farm ID format"})
		return
	}

	medicines, err := h.medicineService.GetAllMedicines(farmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, medicines)
}

// @Summary Get a medicine by ID
// @Description Retrieve a specific medicine by its ID
// @Tags medicines
// @Param id path string true "Medicine ID"
// @Success 200 {object} models.Medicine
// @Failure 400 {object} models.ErrResp
// @Failure 404 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /medicines/{id} [get]
func (h *Handler) GetMedicineByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	medicine, err := h.medicineService.GetMedicineByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if medicine == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "medicine with this ID not found"})
		return
	}

	c.JSON(http.StatusOK, medicine)
}

// @Summary Update an existing medicine
// @Description Update details of a specific medicine
// @Tags medicines
// @Accept application/json
// @Produce application/json
// @Param id path string true "Medicine ID"
// @Param request body models.MedicineReq true "Updated Medicine Details"
// @Success 200 {object} models.MedicineResp
// @Failure 400 {object} models.ErrResp
// @Failure 404 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /medicines/{id} [put]
func (h *Handler) UpdateMedicine(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var medicine models.MedicineWithoutTime
	if err := c.ShouldBindJSON(&medicine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	medicine.ID = id

	if err := h.medicineService.UpdateMedicine(&medicine); err != nil {
		if err == services.ErrNotExist {
			c.JSON(http.StatusNotFound, gin.H{"error": services.ErrNotExist})
		} else if err == services.ErrQuantityLessThanThreshold {
			c.JSON(http.StatusBadRequest, gin.H{"error": services.ErrQuantityLessThanThreshold})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "medicine updated successfully", "medicine": medicine})
}

// @Summary Delete a medicine
// @Description Remove a specific medicine by its ID
// @Tags medicines
// @Param id path string true "Medicine ID"
// @Success 200 {object} models.MessageResp
// @Failure 400 {object} models.ErrResp
// @Failure 404 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /medicines/{id} [delete]
func (h *Handler) DeleteMedicine(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.medicineService.DeleteMedicine(id); err != nil {
		if err == services.ErrNotExist {
			c.JSON(http.StatusNotFound, gin.H{"error": services.ErrNotExist})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medicine deleted successfully"})
}
