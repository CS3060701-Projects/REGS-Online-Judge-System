package models

import (
	"time"
)

// 使用者
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"type:varchar(20);default:'User'"` // 'Admin' 或 'User'
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// 題目
type Problem struct {
	ID           uint   `gorm:"primaryKey"`
	Title        string `gorm:"not null"`
	Description  string `gorm:"type:text"`
	TestcasePath string // 測試資料壓縮檔的路徑
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// 提交紀錄
type Submission struct {
	OperatorID string `gorm:"primaryKey;type:varchar(50)"` // UUID 或隨機字串
	UserID     uint   `gorm:"not null"`
	ProblemID  uint   `gorm:"not null"`
	Status     string `gorm:"type:varchar(20);default:'Pending'"` // Pending, AC, WA, CE, SE, RE, TLE
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// 建立關聯性
	User    User    `gorm:"foreignKey:UserID"`
	Problem Problem `gorm:"foreignKey:ProblemID"`
}
