package handlers

import (
	"errors"
	"net/http"

	"farmish/internal/models"
	"farmish/internal/repository"
	"farmish/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create a new feeding record
// @Tags feeding_records
// @Accept application/json
// @Produce application/json
// @Param request body models.FeedingRecordReq true "Feeding Record request body"
// @Success 201 {object} models.FeedingRecordResp
// @Failure 400 {object} models.ErrResp
// @Failure 404 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /feeding_records [post]
func (h *Handler) CreateFeedingRecord(c *gin.Context) {
	var recordReq models.FeedingRecordWithoutTime
	if err := c.ShouldBindJSON(&recordReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := h.feedingRecordService.CreateFeedingRecord(&recordReq)
	if err != nil {
		if err == services.ErrAnimalNotFound || err == services.ErrFoodNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err == services.ErrInsufficientQuantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Feed record created", "feeding_record": recordReq})
}

// @Summary Get a feeding record by its ID
// @Tags feeding_records
// @Produce application/json
// @Param id path string true "Feeding Record ID"
// @Success 200 {object} models.FeedingRecordDetailed
// @Failure 400 {object} models.ErrResp
// @Failure 404 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /feeding_records/{id} [get]
func (h *Handler) GetFeedingRecordByID(c *gin.Context) {
	id := c.Param("id")
	recordID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feeding record ID"})
		return
	}

	record, err := h.feedingRecordService.GetFeedingRecordByID(recordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch feeding record"})
		return
	}

	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found with this ID"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// @Summary Get all feeding records for a specific animal
// @Tags feeding_records
// @Produce application/json
// @Param animal_id path string true "Animal ID"
// @Success 200 {array} models.FeedingRecordDetailed
// @Failure 400 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /feeding_records/animal/{animal_id} [get]
func (h *Handler) GetFeedingRecordsByAnimalID(c *gin.Context) {
	animalID := c.Param("animal_id")
	parsedAnimalID, err := uuid.Parse(animalID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid animal ID"})
		return
	}

	records, err := h.feedingRecordService.GetFeedingRecordsByAnimalID(parsedAnimalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch feeding records"})
		return
	}

	c.JSON(http.StatusOK, records)
}

// @Summary Update a feeding record by its ID
// @Tags feeding_records
// @Accept application/json
// @Produce application/json
// @Param id path string true "Feeding Record ID"
// @Param input body models.UpdateFeedRecordReq true "Feeding Record Input"
// @Success 200 {object} models.MessageResp
// @Failure 400 {object} models.ErrResp
// @Failure 404 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /feeding_records/{id} [put]
func (h *Handler) UpdateFeedingRecord(c *gin.Context) {
	id := c.Param("id")
	recordID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feeding record ID"})
		return
	}

	var record models.FeedingRecordWithoutTime
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	record.ID = recordID

	err = h.feedingRecordService.UpdateFeedingRecord(&record)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feed record updated successfully"})
}

// @Summary Delete a feeding record by its ID
// @Tags feeding_records
// @Produce application/json
// @Param id path string true "Feeding Record ID"
// @Success 200 {object} models.MessageResp
// @Failure 400 {object} models.ErrResp
// @Failure 404 {object} models.ErrResp
// @Failure 500 {object} models.ErrResp
// @Security BearerAuth
// @Router /feeding_records/{id} [delete]
func (h *Handler) DeleteFeedingRecord(c *gin.Context) {
	id := c.Param("id")
	recordID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feeding record ID"})
		return
	}

	err = h.feedingRecordService.DeleteFeedingRecord(recordID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "feeding record deleted successfully"})
}
