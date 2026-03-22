package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID            string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserName      string         `gorm:"index;not null" json:"username"`
	Title         string         `gorm:"type:varchar(100)" json:"title"`
	Pinned        bool           `gorm:"default:false;not null" json:"pinned"`
	Archived      bool           `gorm:"default:false;not null" json:"archived"`
	LastMessageAt *time.Time     `gorm:"index" json:"last_message_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type SessionInfo struct {
	ID              string     `json:"id,omitempty"`
	SessionID       string     `json:"sessionId"`
	LegacySessionID string     `json:"SessionID,omitempty"`
	Title           string     `json:"name"`
	LegacyTitle     string     `json:"Title,omitempty"`
	Pinned          bool       `json:"pinned"`
	Archived        bool       `json:"archived"`
	LastMessageAt   *time.Time `json:"lastMessageAt,omitempty"`
	UpdatedAt       time.Time  `json:"updatedAt,omitempty"`
}
