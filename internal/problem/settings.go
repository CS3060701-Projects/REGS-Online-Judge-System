package problem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileMapping 代表 settings.yaml 中 replace / public 陣列的元素
type FileMapping struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

// Expected 代表每個 preset 的預期輸出比對方式
type Expected struct {
	Type  string `yaml:"type"`  // "file-content" 或 "string-content"
	Value string `yaml:"value"` // 檔案路徑（相對 problemRoot）或內嵌字串
}

// Preset 代表 settings.yaml 中 presets 陣列的一個測試配置
type Preset struct {
	Score    int           `yaml:"score"`
	Replace  []FileMapping `yaml:"replace"`
	Expected Expected      `yaml:"expected"`
}

// Settings 完整對應 settings.yaml 的結構
type Settings struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Limits      struct {
		TotalTime int `yaml:"totalTime"`
		CpuTime   int `yaml:"cpuTime"`
		Memory    int `yaml:"memory"`
	} `yaml:"limits"`
	Presets []Preset      `yaml:"presets"`
	Public  []FileMapping `yaml:"public"`
}

// LoadSettings 從 problemRoot 載入 settings.yaml
func LoadSettings(problemRoot string) (*Settings, error) {
	path := filepath.Join(problemRoot, "settings.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var settings Settings
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

// MemoryLimitMB 將 KB 轉換為 MB，預設最低 256 MB
func MemoryLimitMB(memoryKB int) int {
	if memoryKB <= 0 {
		return 256
	}
	mb := memoryKB / 1024
	if mb < 1 {
		mb = 1
	}
	return mb
}

// TotalScore 計算所有 preset 的總配分
func (s *Settings) TotalScore() int {
	total := 0
	for _, p := range s.Presets {
		total += p.Score
	}
	return total
}

// LoadExpectedContent 根據 Expected 的 type 載入預期的輸出內容
// problemRoot 是題目根目錄（settings.yaml 所在目錄）
func LoadExpectedContent(expected Expected, problemRoot string) (string, error) {
	switch expected.Type {
	case "file-content":
		path := expected.Value
		if !filepath.IsAbs(path) {
			path = filepath.Join(problemRoot, path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("無法讀取預期輸出檔案 %s: %w", path, err)
		}
		return strings.TrimRight(string(data), "\r\n"), nil
	case "string-content":
		return strings.TrimRight(expected.Value, "\r\n"), nil
	default:
		return "", fmt.Errorf("不支援的 expected type: %s", expected.Type)
	}
}
