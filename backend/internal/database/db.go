package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"regs-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=regs_user password=regs_password dbname=regs_db port=5433 sslmode=disable TimeZone=Asia/Taipei"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連線到資料庫:", err)
	}

	fmt.Println("成功連線到 PostgreSQL!")

	err = DB.AutoMigrate(&models.User{}, &models.Problem{}, &models.Submission{}, &models.JwtBlacklist{})
	if err != nil {
		log.Fatal("資料庫遷移失敗:", err)
	}
	fmt.Println("資料庫遷移完成!")

	syncProblemsFromFolder("testdata")
}
func syncProblemsFromFolder(baseDir string) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		fmt.Printf("無法讀取題目目錄 '%s' (%v)。請確認資料夾存在。\n", baseDir, err)
		return
	}

	existingProblemIDs := make(map[string]bool)
	count := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		problemID := entry.Name()
		problemPath := filepath.Join(baseDir, problemID)

		if _, err := os.Stat(filepath.Join(problemPath, "CMakeLists.txt")); err != nil {
			fmt.Printf("略過題目 '%s'：缺少 CMakeLists.txt\n", problemID)
			continue
		}

		existingProblemIDs[problemID] = true

		var problem models.Problem
		err := DB.Unscoped().Where("id = ?", problemID).First(&problem).Error
		if err == nil {
			updates := map[string]interface{}{
				"testcase_path": problemPath,
				"is_visible":    true,
				"deleted_at":    nil,
			}

			if problem.Title == "" {
				updates["title"] = problemID
			}

			if err := DB.Unscoped().Model(&problem).Updates(updates).Error; err != nil {
				fmt.Printf("更新題目 '%s' 失敗: %v\n", problemID, err)
				continue
			}
		} else if err == gorm.ErrRecordNotFound {
			problem = models.Problem{
				ID:           problemID,
				Title:        problemID,
				TestcasePath: problemPath,
				IsVisible:    true,
			}

			if err := DB.Create(&problem).Error; err != nil {
				fmt.Printf("建立題目 '%s' 失敗: %v\n", problemID, err)
				continue
			}
		} else {
			fmt.Printf("查詢題目 '%s' 失敗: %v\n", problemID, err)
			continue
		}

		count++
	}

	var dbProblems []models.Problem
	if err := DB.Find(&dbProblems).Error; err != nil {
		fmt.Printf("讀取資料庫題目列表失敗: %v\n", err)
		return
	}

	hiddenCount := 0
	for _, problem := range dbProblems {
		if !existingProblemIDs[problem.ID] {
			if err := DB.Model(&models.Problem{}).
				Where("id = ?", problem.ID).
				Update("is_visible", false).Error; err != nil {
				fmt.Printf("隱藏已不存在的題目 '%s' 失敗: %v\n", problem.ID, err)
				continue
			}

			hiddenCount++
		}
	}

	fmt.Printf("題目初始化完成！從 '%s' 載入 %d 題，隱藏 %d 題不存在的題目。\n", baseDir, count, hiddenCount)
}
