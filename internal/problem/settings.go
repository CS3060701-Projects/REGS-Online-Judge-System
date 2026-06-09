package problem

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Limits      struct {
		TotalTime int `yaml:"totalTime"`
		CpuTime   int `yaml:"cpuTime"`
		Memory    int `yaml:"memory"`
	} `yaml:"limits"`
}

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
