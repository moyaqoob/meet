package models

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"size:100;not null"`
	Email        string `gorm:"uniqueIndex;size:255;not null"`
	Password     string `gorm:"not null"` // hashed
	AvatarURL    *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Meeting      []Meeting            `gorm:"foreignKey:HostID"`
	Participants []MeetingParticipant `gorm:"foreignKey:UserID"`
	Message      []Message            `gorm:"foreignKey:SenderID"`
}

type Meeting struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	HostID    uint   `gorm:"not null;index"`
	Title     string `gorm:"size:300"`
	MeetingID string   `gorm:"index;size:10"`
	ScheduleAt time.Time
	IsActive   bool
	CreatedAt  time.Time

	Host User `gorm:"foreignKey:HostID"`

	Participants []MeetingParticipant `gorm:"foreignKey:MeetingID"`
	Messages     []Message            `gorm:"foreignKey:MeetingID"`
}

type MeetingParticipant struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	MeetingID uint   `gorm:"index;not null"`
	UserID    uint   `gorm:"index;not null"`
	Role      string `gorm:"size:50"` //
	JoinedAt  time.Time
	LeftAt    *time.Time
}

type Message struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	MeetingID uint   `gorm:"index;not null"`
	SenderID  uint   `gorm:"index;not null"`
	Content   string `gorm:"type:text;not null"`
	SentAt    time.Time
}
