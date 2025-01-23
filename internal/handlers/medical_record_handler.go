package handlers

import (
	"errors"
	"farmish/internal/models"
	"farmish/internal/repository"
	"farmish/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create a new medical record
// @Tags medical_records
// @Accept application/json
// @Produce application/json
// @Param request body models.MedicalRecordReq true "Medical Record"
// @Success 201 {object} models.MedicalRecordResp "Created successfully"
// @Failure 400 {object} models.ErrResp "Invalid input"
// @Failure 404 {object} models.ErrResp "Animal or Medicine Not Found"
// @Failure 500 {object} models.ErrResp "Internal server error"
// @Security BearerAuth
// @Router /medical_records [post]
func (h *Handler) CreateMedicalRecord(c *gin.Context) {
	var record models.MedicalRecordWithoutTime
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	err := h.medicalRecordService.CreateMedicalRecord(&record)
	if err != nil {
		if err == services.ErrAnimalNotFound || err == services.ErrMedicineNotExist {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err == services.ErrInsufficientQuantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created successfully", "medical_record": record})
}

// @Summary Get a medical record by ID
// @Tags medical_records
// @Produce application/json
// @Param id path string true "Medical Record ID"
// @Success 200 {object} models.MedicalRecordDetailed
// @Failure 400 {object} models.ErrResp "Invalid record ID"
// @Failure 404 {object} models.ErrResp "Not found"
// @Failure 500 {object} models.ErrResp "Internal server error"
// @Security BearerAuth
// @Router /medical_records/{id} [get]
func (h *Handler) GetMedicalRecordByID(c *gin.Context) {
	id := c.Param("id")
	recordID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid medical record ID"})
		return
	}

	record, err := h.medicalRecordService.GetMedicalRecordByID(recordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// @Summary Get medical records by Animal ID
// @Tags medical_records
// @Produce application/json
// @Param animal_id path string true "Animal ID"
// @Success 200 {array} models.MedicalRecordDetailed
// @Failure 400 {object} models.ErrResp "Invalid animal ID"
// @Failure 500 {object} models.ErrResp "Internal server error"
// @Security BearerAuth
// @Router /medical_records/animals/{animal_id} [get]
func (h *Handler) GetMedicalRecordsByAnimalID(c *gin.Context) {
	id := c.Param("animal_id")
	animalID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid animal record ID"})
		return
	}

	records, err := h.medicalRecordService.GetMedicalRecordsByAnimalID(animalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// @Summary Update a medical record
// @Tags medical_records
// @Accept  application/json
// @Produce application/json
// @Param id path string true "Medical Record ID"
// @Param request body models.UpdateMedicalRecordReq true "Medical Record"
// @Success 200 {object} models.MessageResp "Updated successfully"
// @Failure 400 {object} models.ErrResp "Invalid input"
// @Failure 404 {object} models.ErrResp "Not Found"
// @Failure 500 {object} models.ErrResp "Internal server error"
// @Security BearerAuth
// @Router /medical_records/{id} [put]
func (h *Handler) UpdateMedicalRecord(c *gin.Context) {
	id := c.Param("id")
	recordID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feeding record ID"})
		return
	}

	var record models.MedicalRecordWithoutTime
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	record.ID = recordID

	err = h.medicalRecordService.UpdateMedicalRecord(&record)
	if err != nil {
		if errors.Is(err, repository.ErrMedicalRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrMedicalRecordNotFound})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

// @Summary Delete a medical record by ID
// @Tags medical_records
// @Produce application/json
// @Param id path string true "Medical Record ID"
// @Success 200 {object} models.MessageResp "Deleted successfully"
// @Failure 400 {object} models.ErrResp "Invalid record ID"
// @Failure 404 {object} models.ErrResp "Not found"
// @Failure 500 {object} models.ErrResp "Internal server error"
// @Security BearerAuth
// @Router /medical_records/{id} [delete]
func (h *Handler) DeleteMedicalRecord(c *gin.Context) {
	id := c.Param("id")
	recordID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid medical record ID"})
		return
	}

	err = h.medicalRecordService.DeleteMedicalRecord(recordID)
	if err != nil {
		if errors.Is(err, repository.ErrMedicalRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": repository.ErrMedicalRecordNotFound})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}
