package main

import (
	"log"
	"mynote/internal/database"
	"mynote/internal/handlers"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	// ハンドラの初期化
	handler := handlers.NewHandler(db.DB)

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
