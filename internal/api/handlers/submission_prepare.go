package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"regs-backend/internal/models"
	"regs-backend/pkg/utils"
)

var officialInjectFiles = []string{
	"entrypoint.cpp",
	"test.h",
}

func problemRoot(problem models.Problem) string {
	if problem.TestcasePath != "" {
		return problem.TestcasePath
	}
	return filepath.Join("testdata", problem.ID)
}

func hasCMakeLists(workspace string) bool {
	_, err := os.Stat(filepath.Join(workspace, "CMakeLists.txt"))
	return err == nil
}

func ensureCMakeProject(workspace string, problem models.Problem) error {
	if hasCMakeLists(workspace) {
		return nil
	}

	root := problemRoot(problem)
	candidates := []string{
		filepath.Join(root, "template", "CMakeLists.txt"),
		filepath.Join(root, "solution", "CMakeLists.txt"),
		filepath.Join(root, "CMakeLists.txt"),
	}

	for _, src := range candidates {
		if _, err := os.Stat(src); err != nil {
			continue
		}
		if err := utils.CopyFile(src, filepath.Join(workspace, "CMakeLists.txt")); err != nil {
			return err
		}
		break
	}

	if !hasCMakeLists(workspace) {
		return fmt.Errorf("上傳的專案根目錄缺少 CMakeLists.txt")
	}

	specSrc := filepath.Join(root, "spec")
	if info, err := os.Stat(specSrc); err == nil && info.IsDir() {
		specDst := filepath.Join(workspace, "spec")
		_ = os.RemoveAll(specDst)
		if err := utils.CopyDir(specSrc, specDst); err != nil {
			return fmt.Errorf("複製 spec 失敗: %w", err)
		}
	}

	return nil
}

func injectOfficialProblemFiles(workspace string, problem models.Problem) error {
	root := problemRoot(problem)

	for _, fileName := range officialInjectFiles {
		officialPath := filepath.Join(root, "solution", fileName)
		if _, err := os.Stat(officialPath); err != nil {
			return fmt.Errorf("題目官方 %s 不存在: %s", fileName, officialPath)
		}

		if err := removeFilesByName(workspace, fileName); err != nil {
			return err
		}

		content, err := os.ReadFile(officialPath)
		if err != nil {
			return err
		}

		if err := os.WriteFile(filepath.Join(workspace, fileName), content, 0644); err != nil {
			return err
		}
	}

	return nil
}

func removeFilesByName(workspace, fileName string) error {
	return filepath.WalkDir(workspace, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.EqualFold(d.Name(), fileName) {
			return os.Remove(path)
		}
		return nil
	})
}

func mirrorConfigLog(workspace string) {
	configurePath := filepath.Join(workspace, "configure.log")
	configPath := filepath.Join(workspace, "config.log")
	data, err := os.ReadFile(configurePath)
	if err != nil {
		return
	}
	_ = os.WriteFile(configPath, data, 0644)
}
