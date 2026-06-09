package models

import (
	"time"

	"gorm.io/gorm"
)

const JUDGER_IMAGE = "regs-judger"

// PresetResult 記錄單一 preset 的評測結果
type PresetResult struct {
	Index    int     `json:"index"`
	Status   string  `json:"status"`    // AC, WA, CE, RE, TLE, SE
	Score    int     `json:"score"`     // 該 preset 配分
	Earned   int     `json:"earned"`    // 實得分數 (AC=score, 其他=0)
	PeakTime float64 `json:"peak_time"`
}

type JudgeResult struct {
	Status        string         `json:"status"`
	PeakTime      float64        `json:"peak_time"`
	PeakMemory    int64          `json:"peak_memory"`
	TotalScore    int            `json:"total_score"`
	EarnedScore   int            `json:"earned_score"`
	PresetResults []PresetResult `json:"preset_results,omitempty"`
}

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null" json:"-"`
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

	Score      int `json:"score"`       // 實際得分
	TotalScore int `json:"total_score"` // 該題滿分
}

type JwtBlacklist struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"type:text;index"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
}
