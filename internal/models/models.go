package models

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"type:varchar(20);default:'User'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Problem struct {
	ID           string `gorm:"primaryKey;type:varchar(50)"`
	Title        string `gorm:"not null"`
	Description  string `gorm:"type:text"`
	TestcasePath string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Submission struct {
	OperatorID string `gorm:"primaryKey;type:varchar(50)"`
	UserID     uint   `gorm:"not null"`
	ProblemID  string `gorm:"not null"`
	Status     string `gorm:"type:varchar(20);default:'Pending'"` // Pending, AC, WA, CE, SE, RE, TLE
	CreatedAt  time.Time
	UpdatedAt  time.Time

	User    User    `gorm:"foreignKey:UserID"`
	Problem Problem `gorm:"foreignKey:ProblemID"`
}
