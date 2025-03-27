package room

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"polling-app/internal/models"
	"polling-app/pkg/database"
)

type CreateRoomRequest struct {
	Name string `json:"name" binding:"required"`
}

type JoinRoomRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
}

func CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	// Create new room
	room := models.Room{
		Name:   req.Name,
		HostID: currentUser.ID,
	}

	if err := database.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	// Add host as first participant
	if err := database.DB.Model(&room).Association("Participants").Append(&currentUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add host to participants"})
		return
	}

	c.JSON(http.StatusCreated, room)
}

func GetRoom(c *gin.Context) {
	roomID := c.Param("id")

	var room models.Room
	if err := database.DB.Preload("Host").Preload("Participants").First(&room, "id = ?", roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}

func JoinRoom(c *gin.Context) {
	var req JoinRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	var room models.Room
	if err := database.DB.First(&room, "invite_code = ?", req.InviteCode).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	if !room.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room is not active"})
		return
	}

	// Check if user is already a participant
	var count int64
	database.DB.Model(&room).Where("user_id = ?", currentUser.ID).Association("Participants").Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already a participant in this room"})
		return
	}

	// Add user to participants
	if err := database.DB.Model(&room).Association("Participants").Append(&currentUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join room"})
		return
	}

	c.JSON(http.StatusOK, room)
} 