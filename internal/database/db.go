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

	// 執行資料表結構自動遷移
	err = DB.AutoMigrate(&models.User{}, &models.Problem{}, &models.Submission{})
	if err != nil {
		log.Fatal("資料庫遷移失敗:", err)
	}
	fmt.Println("資料庫遷移完成!")

	syncProblemsFromFolder("test_data")
}

func syncProblemsFromFolder(baseDir string) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		fmt.Printf("無法讀取題目目錄 '%s' (%v)。請確認資料夾存在。\n", baseDir, err)
		return
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			problemID := entry.Name()

			problem := models.Problem{
				ID:           problemID,
				Title:        "自動匯入: " + problemID,              // 暫時用 ID 當標題
				TestcasePath: filepath.Join(baseDir, problemID), // 把測資路徑記錄下來
			}

			DB.FirstOrCreate(&problem, models.Problem{ID: problemID})
			count++
		}
	}

	fmt.Printf("題目初始化完成！共從 '%s' 載入 %d 題。\n", baseDir, count)
}
