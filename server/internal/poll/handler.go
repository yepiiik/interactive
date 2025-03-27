package poll

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"polling-app/internal/models"
	"polling-app/internal/websocket"
	"polling-app/pkg/database"
)

type CreatePollRequest struct {
	RoomID    string   `json:"room_id" binding:"required"`
	Question  string   `json:"question" binding:"required"`
	Options   []string `json:"options" binding:"required,min=2,max=4"`
	Duration  int      `json:"duration" binding:"required,min=5,max=300"` // Duration in seconds
	CorrectID uint     `json:"correct_id" binding:"required"`
}

type VoteRequest struct {
	OptionID  uint    `json:"option_id" binding:"required"`
	TimeTaken float64 `json:"time_taken" binding:"required"` // Time taken to answer in seconds
}

func CreatePoll(c *gin.Context) {
	var req CreatePollRequest
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

	// Verify user is the host of the room
	var room models.Room
	if err := database.DB.First(&room, "id = ?", req.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	if room.HostID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the host can create polls"})
		return
	}

	// Create poll
	poll := models.Poll{
		RoomID:   req.RoomID,
		Question: req.Question,
		Duration: req.Duration,
	}

	if err := database.DB.Create(&poll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create poll"})
		return
	}

	// Create options
	for i, optionText := range req.Options {
		option := models.Option{
			PollID:    poll.ID,
			Text:      optionText,
			IsCorrect: uint(i+1) == req.CorrectID,
		}
		if err := database.DB.Create(&option).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create options"})
			return
		}
	}

	// Start the poll
	poll.StartPoll()
	if err := database.DB.Save(&poll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start poll"})
		return
	}

	// Broadcast poll start to all clients in the room
	websocket.BroadcastToRoom(req.RoomID, "start_poll", poll)

	// Set timer to end poll
	go func() {
		time.Sleep(time.Duration(req.Duration) * time.Second)
		endPoll(poll.ID)
	}()

	c.JSON(http.StatusCreated, poll)
}

func Vote(c *gin.Context) {
	pollID := c.Param("id")
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(models.User)

	// Verify poll exists and is active
	var poll models.Poll
	if err := database.DB.First(&poll, "id = ?", pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	if !poll.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll is not active"})
		return
	}

	// Check if user has already voted
	var existingVote models.Vote
	if err := database.DB.Where("poll_id = ? AND user_id = ?", pollID, currentUser.ID).First(&existingVote).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already voted"})
		return
	}

	// Create vote
	vote := models.Vote{
		UserID:    currentUser.ID,
		PollID:    poll.ID,
		OptionID:  req.OptionID,
		TimeTaken: req.TimeTaken,
	}

	if err := database.DB.Create(&vote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vote"})
		return
	}

	// Broadcast vote to all clients
	websocket.BroadcastToRoom(poll.RoomID, "vote", vote)

	c.JSON(http.StatusOK, vote)
}

func GetResults(c *gin.Context) {
	pollID := c.Param("id")

	var poll models.Poll
	if err := database.DB.Preload("Options.Votes").First(&poll, "id = ?", pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	// Calculate results
	type OptionResult struct {
		OptionID   uint    `json:"option_id"`
		Text       string  `json:"text"`
		VoteCount  int     `json:"vote_count"`
		IsCorrect  bool    `json:"is_correct"`
		Percentage float64 `json:"percentage"`
	}

	var results []OptionResult
	totalVotes := 0

	// Count total votes
	for _, option := range poll.Options {
		totalVotes += len(option.Votes)
	}

	// Calculate results for each option
	for _, option := range poll.Options {
		voteCount := len(option.Votes)
		percentage := 0.0
		if totalVotes > 0 {
			percentage = float64(voteCount) / float64(totalVotes) * 100
		}

		results = append(results, OptionResult{
			OptionID:   option.ID,
			Text:      option.Text,
			VoteCount: voteCount,
			IsCorrect: option.IsCorrect,
			Percentage: percentage,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"poll":    poll,
		"results": results,
	})
}

func endPoll(pollID uint) {
	var poll models.Poll
	if err := database.DB.First(&poll, "id = ?", pollID).Error; err != nil {
		return
	}

	poll.IsActive = false
	if err := database.DB.Save(&poll).Error; err != nil {
		return
	}

	// Broadcast poll end to all clients
	websocket.BroadcastToRoom(poll.RoomID, "end_poll", poll)
} 