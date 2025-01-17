package handlers

import (
	"farmish/internal/models"
	"farmish/internal/repository"
	"farmish/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		User login
// @Description	Authenticate a user using their email and password.
// @Tags			auth
// @Accept			application/json
// @Produce		application/json
// @Param			request	body		models.LoginRequest	true	"Login credentials"
// @Success		200		{object}	models.LoginResponse	"Successful login response with token and user ID"
// @Failure		400		{object}	models.ErrResp		"Invalid input format or missing fields"
// @Failure		401		{object}	models.ErrResp		"Invalid email or password"
// @Failure		500		{object}	models.ErrResp		"Internal server error"
// @Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	credentials := models.LoginRequest{}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := h.userService.Login(&credentials)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": resp.Token, "user_id": resp.ID})
}

// @Summary User sign-up
// @Description Creates a new user account.
//
//	@Tags			auth
//
// @Accept application/json
// @Produce application/json
// @Param request body models.SignUpRequest true "Sign-up details"
// @Success 201 {object} models.SignUpResponse
// @Failure 400 {object} models.ErrResp
//
//	@Failure		409		{object}	models.ErrResp		"Email exist"
//	@Failure		500		{object}	models.ErrResp		"Internal server error"
//
// @Router /auth/signup [post]
func (h *Handler) SignUp(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.userService.SignUp(&user)
	if err != nil {
		if err == repository.ErrEmailAlreadyInUse {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   token,
		"user_id": user.ID,
	})
}
