package handler

import (
	"net/http"
	"strconv"

	"github.com/JunBumHan/copilot-agent-test/internal/domain"
	"github.com/JunBumHan/copilot-agent-test/internal/usecase"
	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler creates a new UserHandler with the given use case
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetUser handles GET requests to retrieve a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	user, err := h.userUseCase.GetByID(uint(id))
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers handles GET requests to retrieve all users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userUseCase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// CreateUser handles POST requests to create a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.userUseCase.Create(&user); err != nil {
		switch err {
		case usecase.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser handles PUT requests to update an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user.ID = uint(id)

	if err := h.userUseCase.Update(&user); err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case usecase.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE requests to remove a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.userUseCase.Delete(uint(id)); err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// RegisterRoutes registers the user routes with the given router group
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("", h.GetAllUsers)
		users.GET("/:id", h.GetUser)
		users.POST("", h.CreateUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}