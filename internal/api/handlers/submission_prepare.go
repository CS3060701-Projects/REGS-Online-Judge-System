package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"regs-backend/internal/models"
	"regs-backend/internal/problem"
	"regs-backend/pkg/utils"
)

func problemRoot(prob models.Problem) string {
	if prob.TestcasePath != "" {
		return prob.TestcasePath
	}
	return filepath.Join("testdata", prob.ID)
}

func hasCMakeLists(workspace string) bool {
	_, err := os.Stat(filepath.Join(workspace, "CMakeLists.txt"))
	return err == nil
}

func ensureCMakeProject(workspace string, prob models.Problem) error {
	if hasCMakeLists(workspace) {
		return nil
	}

	root := problemRoot(prob)

	// 嘗試從 settings.yaml 的 public 欄位找 CMakeLists.txt
	settings, err := problem.LoadSettings(root)
	if err == nil && len(settings.Public) > 0 {
		for _, p := range settings.Public {
			if p.Target == "CMakeLists.txt" {
				srcPath := p.Source
				if !filepath.IsAbs(srcPath) {
					srcPath = filepath.Join(root, srcPath)
				}
				if _, err := os.Stat(srcPath); err == nil {
					if err := utils.CopyFile(srcPath, filepath.Join(workspace, "CMakeLists.txt")); err != nil {
						return err
					}
					break
				}
			}
		}
	}

	// 若 settings.yaml 沒提供，fallback 到舊邏輯
	if !hasCMakeLists(workspace) {
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

// injectOfficialProblemFiles 使用 settings.yaml 的 replace 規則已移至 judge 模組中逐 preset 處理。
// 此函式保留為向後相容（用於沒有 settings.yaml 的題目）。
func injectOfficialProblemFiles(workspace string, prob models.Problem) error {
	root := problemRoot(prob)

	// 若存在 settings.yaml，注入在 judge 流程中逐 preset 處理，此處不需要做
	if _, err := problem.LoadSettings(root); err == nil {
		fmt.Printf("[injectOfficialProblemFiles] 偵測到 settings.yaml，跳過全域注入（將由 judge 逐 preset 處理）\n")
		return nil
	}

	// Fallback: 沒有 settings.yaml 的舊邏輯
	officialInjectFiles := []string{
		"entrypoint.cpp",
		"test.h",
	}

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
		if d.Name() == fileName {
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
