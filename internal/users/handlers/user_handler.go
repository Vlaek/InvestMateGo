package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"invest-mate/internal/users/models/domain"
	"invest-mate/internal/users/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("/register", h.Register)
		users.POST("/login", h.Login)

		// Защищенные маршруты (потребуют middleware авторизации)
		users.GET("/profile", h.GetProfile)
		users.PUT("/profile", h.UpdateProfile)
		// users.DELETE("/profile", h.DeactivateProfile)
	}

	// Admin routes
	admin := router.Group("/admin/users")
	{
		admin.GET("/", h.ListUsers)
		admin.GET("/:id", h.GetUserByID)
		admin.PUT("/:id", h.UpdateUser)
		// admin.DELETE("/:id", h.DeactivateUser)
	}
}

// Register регистрирует нового пользователя
func (h *UserHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest

	// Парсим JSON тело запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Вызываем сервис
	userResponse, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		errorMsg := err.Error()

		// Обработка специфических ошибок
		if err.Error() == "email already exists" ||
			err.Error() == "username already taken" ||
			err.Error() == "password must be at least 8 characters" {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{"error": errorMsg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    userResponse,
	})
}

// Login выполняет вход пользователя
func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userResponse, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "invalid credentials" ||
			err.Error() == "account is deactivated" {
			status = http.StatusUnauthorized
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// TODO: Генерация JWT токена
	// token, err := generateJWT(userResponse.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    userResponse,
		// "token":   token,
	})
}

// GetProfile возвращает профиль текущего пользователя
func (h *UserHandler) GetProfile(c *gin.Context) {
	// TODO: Получать ID пользователя из JWT токена
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userResponse, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userResponse})
}

// UpdateProfile обновляет профиль пользователя
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var updates domain.User
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userResponse, err := h.userService.UpdateUser(c.Request.Context(), userID, &updates)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email already exists" ||
			err.Error() == "username already taken" {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    userResponse,
	})
}

// ListUsers возвращает список пользователей (только для админов)
func (h *UserHandler) ListUsers(c *gin.Context) {
	// TODO: Проверка прав администратора

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, err := h.userService.ListUsers(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
		},
	})
}

// GetUserByID возвращает пользователя по ID (админ)
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	userResponse, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userResponse})
}

// UpdateUser обновляет пользователя (админ)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var updates domain.User
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userResponse, err := h.userService.UpdateUser(c.Request.Context(), userID, &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    userResponse,
	})
}
