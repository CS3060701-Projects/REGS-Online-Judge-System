package database

import (
	"fmt"
	"log"

	"regs-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// 這裡的帳密對應剛剛 docker-compose.yml 內的設定
	dsn := "host=localhost user=regs_user password=regs_password dbname=regs_db port=5433 sslmode=disable TimeZone=Asia/Taipei"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連線到資料庫:", err)
	}

	fmt.Println("成功連線到 PostgreSQL!")

	// 自動建立或更新資料表結構
	err = DB.AutoMigrate(&models.User{}, &models.Problem{}, &models.Submission{})
	if err != nil {
		log.Fatal("資料庫遷移失敗:", err)
	}
	fmt.Println("資料庫遷移完成!")
}
