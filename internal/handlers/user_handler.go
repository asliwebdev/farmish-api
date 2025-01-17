package handlers

import (
	"farmish/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary		Get a user by ID
// @Description	Retrieve user details by their UUID.
// @Tags			users
// @Produce		application/json
// @Param			id		path		string			true	"User ID (UUID)"
// @Success		200		{object}	models.User		"User details"
// @Failure		400		{object}	models.ErrResp	"Invalid user ID format"
// @Failure		404		{object}	models.ErrResp	"User not found"
// @Failure		500		{object}	models.ErrResp	"Internal server error"
// @Security		BearerAuth
// @Router			/users/{id} [get]
func (h *Handler) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is not in the UUID format"})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary		Get all users
// @Description	Retrieve a list of all registered users.
// @Tags			users
// @Produce		application/json
// @Success		200		{array}		models.User		"List of users"
// @Failure		500		{object}	models.ErrResp	"Internal server error"
// @Security		BearerAuth
// @Router			/users [get]
func (h *Handler) GetAllUsers(ctx *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// @Summary		Update a user
// @Description	Update details of a user by their UUID.
// @Tags			users
// @Accept			application/json
// @Produce		application/json
// @Param			id		path		string			true	"User ID (UUID)"
// @Param			request	body		models.UpdateUserSwag	true	"User update payload, provide password only if it is updated"
// @Success		200		{object}	models.UpdateUserResp			"User updated successfully"
// @Failure		400		{object}	models.ErrResp	"Invalid input or user ID format"
// @Failure		500		{object}	models.ErrResp	"Internal server error"
// @Security		BearerAuth
// @Router			/users/{id} [put]
func (h *Handler) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is not in the UUID format"})
		return
	}

	user := models.UpdateUser{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	user.ID = userId

	if user.Password != "" && len(user.Password) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password length should not be less than 6"})
		return
	}

	err = h.userService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

// @Summary		Delete a user
// @Description	Delete a user by their UUID.
// @Tags			users
// @Produce		application/json
// @Param			id		path		string			true	"User ID (UUID)"
// @Success		200		{object}	models.MessageResp			"User deleted successfully"
// @Failure		400		{object}	models.ErrResp	"Invalid user ID format"
// @Failure		404		{object}	models.ErrResp	"User not found"
// @Failure		500		{object}	models.ErrResp	"Internal server error"
// @Security		BearerAuth
// @Router			/users/{id} [delete]
func (h *Handler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is not in the UUID format"})
		return
	}

	err = h.userService.DeleteUser(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
