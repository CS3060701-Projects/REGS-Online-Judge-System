package models

import (
	"time"

	"gorm.io/gorm"
)

const JUDGER_IMAGE = "regs-judger"

type JudgeResult struct {
	Status     string
	PeakTime   float64
	PeakMemory int64
}

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"type:varchar(20);default:'User'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Problem struct {
	ID          string `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Title       string `gorm:"not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`

	TimeLimit   int `gorm:"default:1000" json:"time_limit"`
	MemoryLimit int `gorm:"default:256" json:"memory_limit"`

	TestcasePath string `json:"testcase_path"`
	IsVisible    bool   `gorm:"default:true" json:"is_visible"`

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Submissions []Submission   `gorm:"foreignKey:ProblemID" json:"-"`
}

type Submission struct {
	OperatorID string `gorm:"primaryKey;type:varchar(50)"`
	UserID     uint   `gorm:"not null"`
	ProblemID  string `gorm:"not null"`
	Status     string `gorm:"type:varchar(20);default:'Pending'"` // Pending, AC, WA, CE, SE, RE, TLE
	CreatedAt  time.Time
	UpdatedAt  time.Time

	User      User    `gorm:"foreignKey:UserID"`
	Problem   Problem `gorm:"foreignKey:ProblemID"`
	RunTime   int     `json:"run_time"`
	RunMemory int64   `json:"run_memory"`
}

type JwtBlacklist struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"type:text;index"` // 儲存整個 Token 字串
	ExpiresAt time.Time `gorm:"index"`           // Token 原本的過期時間
	CreatedAt time.Time
}
