package main

import (
	"fmt"
	"log"
	"mynote/internal/handlers"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Data Source Name with charset and parseTime
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// ハンドラの初期化
	handler := handlers.NewHandler(db)

	// ルートの定義
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/register-form", handler.RegisterForm)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/select", handler.Select)
	http.HandleFunc("/update", handler.Update)
	http.HandleFunc("/delete", handler.Delete)

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
