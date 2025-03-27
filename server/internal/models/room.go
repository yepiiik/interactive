package models

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	HostID      uint      `json:"host_id" gorm:"not null"`
	Host        User      `json:"host" gorm:"foreignKey:HostID"`
	InviteCode  string    `json:"invite_code" gorm:"unique;not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Participants []User   `json:"participants" gorm:"many2many:room_participants;"`
}

// BeforeCreate generates a unique invite code before creating the room
func (r *Room) BeforeCreate() error {
	r.ID = uuid.New().String()
	r.InviteCode = generateInviteCode()
	return nil
}

// generateInviteCode creates a random 6-character invite code
func generateInviteCode() string {
	// Generate a random 6-character code using alphanumeric characters
	code := uuid.New().String()[:6]
	return code
} 