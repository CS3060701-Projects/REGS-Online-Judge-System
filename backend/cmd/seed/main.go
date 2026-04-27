package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regs-backend/internal/database"
	"regs-backend/internal/models"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("Connecting to the database...")
	database.Connect()

	fmt.Println("No admin user found. Let's create one.")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter admin username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter admin password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		log.Fatal("Username and password cannot be empty.")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	adminUser := models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         "Admin",
	}

	if err := database.DB.Create(&adminUser).Error; err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Printf("Admin user '%s' created successfully!\n", username)
}
