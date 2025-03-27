package models

import (
	"time"
)

type Poll struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	RoomID      string    `json:"room_id" gorm:"not null"`
	Room        Room      `json:"room" gorm:"foreignKey:RoomID"`
	Question    string    `json:"question" gorm:"not null"`
	Options     []Option  `json:"options" gorm:"foreignKey:PollID"`
	Duration    int       `json:"duration" gorm:"not null"` // Duration in seconds
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsActive    bool      `json:"is_active" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Option struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PollID    uint      `json:"poll_id" gorm:"not null"`
	Text      string    `json:"text" gorm:"not null"`
	IsCorrect bool      `json:"is_correct" gorm:"default:false"`
	Votes     []Vote    `json:"votes" gorm:"foreignKey:OptionID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Vote struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	PollID    uint      `json:"poll_id" gorm:"not null"`
	Poll      Poll      `json:"poll" gorm:"foreignKey:PollID"`
	OptionID  uint      `json:"option_id" gorm:"not null"`
	Option    Option    `json:"option" gorm:"foreignKey:OptionID"`
	TimeTaken float64   `json:"time_taken" gorm:"not null"` // Time taken to answer in seconds
	CreatedAt time.Time `json:"created_at"`
}

// StartPoll activates the poll and sets the start and end times
func (p *Poll) StartPoll() {
	p.IsActive = true
	p.StartTime = time.Now()
	p.EndTime = p.StartTime.Add(time.Duration(p.Duration) * time.Second)
}

// IsExpired checks if the poll has ended
func (p *Poll) IsExpired() bool {
	return time.Now().After(p.EndTime)
}

// GetTimeRemaining returns the remaining time in seconds
func (p *Poll) GetTimeRemaining() float64 {
	if !p.IsActive {
		return 0
	}
	remaining := p.EndTime.Sub(time.Now()).Seconds()
	if remaining < 0 {
		return 0
	}
	return remaining
} 