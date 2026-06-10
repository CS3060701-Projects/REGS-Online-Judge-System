package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "server":
		runServer()
	case "reset-db":
		resetDB()
	case "seed-admin":
		seedAdmin()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("REGS Task Runner (Cross-Platform)")
	fmt.Println("Usage: go run cmd/task/main.go <command>")
	fmt.Println("\nCommands:")
	fmt.Println("  server       - Ensure Docker is running, build, and start the server")
	fmt.Println("  reset-db     - Reset the database (DELETES ALL DATA)")
	fmt.Println("  seed-admin   - Create the default admin user")
}

func ensureDocker() {
	fmt.Println("[INFO] Ensuring PostgreSQL container is running...")
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("[ERROR] Failed to start Docker. Is Docker service running?")
		os.Exit(1)
	}
	fmt.Println("[INFO] Waiting for database to be ready...")
	time.Sleep(4 * time.Second)
}

func runServer() {
	fmt.Println("=========================================")
	fmt.Println("        REGS 評測系統 - 伺服器啟動")
	fmt.Println("=========================================")
	fmt.Println()

	ensureDocker()

	fmt.Println("\n[提示] Docker 容器已在運行中。")
	fmt.Println("\n[2/3] 正在編譯最新程式碼...")

	binaryName := "server"
	if runtime.GOOS == "windows" {
		binaryName = "server.exe"
	}

	os.Remove(binaryName)

	buildCmd := exec.Command("go", "build", "-o", binaryName, "./cmd/server")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		fmt.Println("\n[錯誤] 編譯失敗！請檢查上方的錯誤訊息。")
		os.Exit(1)
	}

	fmt.Println("\n[3/3] 編譯成功！準備啟動伺服器...")
	fmt.Println("-----------------------------------------")
	fmt.Println("提示：若要關閉伺服器，請按 Ctrl+C")
	fmt.Println("-----------------------------------------")
	fmt.Println()

	runCmd := exec.Command("." + string(filepath.Separator) + binaryName)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Stdin = os.Stdin
	if err := runCmd.Run(); err != nil {
		fmt.Printf("\n[錯誤] 伺服器異常終止: %v\n", err)
		os.Exit(1)
	}
}

func seedAdmin() {
	fmt.Println("=========================================")
	fmt.Println("      REGS System - Create Admin User")
	fmt.Println("=========================================")
	fmt.Println()

	ensureDocker()

	seedCmd := exec.Command("go", "run", "./cmd/seed")
	seedCmd.Stdout = os.Stdout
	seedCmd.Stderr = os.Stderr
	seedCmd.Stdin = os.Stdin
	if err := seedCmd.Run(); err != nil {
		fmt.Println("\n[錯誤] 建立管理員失敗。")
		os.Exit(1)
	}
	fmt.Println("\nOperation finished.")
}

func resetDB() {
	fmt.Println("=========================================")
	fmt.Println("     REGS System - Reset Database")
	fmt.Println("=========================================")
	fmt.Println()
	fmt.Println("WARNING: This will completely and IRREVERSIBLY delete all data")
	fmt.Println("from the PostgreSQL database, including all users, problems,")
	fmt.Println("and submissions.")
	fmt.Println()

	fmt.Print("Are you sure you want to continue? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" {
		fmt.Println("Operation cancelled.")
		return
	}

	fmt.Println("\nStopping database container and removing data volume...")
	cmd := exec.Command("docker", "compose", "down", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("[錯誤] 重設資料庫失敗。")
		os.Exit(1)
	}

	fmt.Println("\nDatabase has been successfully reset.")
	fmt.Println("You can now restart it using: go run cmd/task/main.go server")
}
