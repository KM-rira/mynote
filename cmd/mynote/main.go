package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Note struct {
	ID            int
	Title         string
	Content       string
	Category      string
	ImportantFlag bool
	UpdatedAt     time.Time
}

func main() {
	// ポート8081でサーバーを起動します。
	fmt.Println("Server is running on http://localhost:8081")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("access /")
		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("access /index")
		// データの取得
		rows, err := db.Query("SELECT * FROM note")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to execute query", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// 結果を表示
		var notes []Note
		for rows.Next() {
			var note Note
			var updatedAtStr []uint8
			if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.Category, &note.ImportantFlag, &updatedAtStr); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				log.Println("Failed to scan row:", err)
				return
			}
			// updatedAtStr を time.Time に変換
			updatedAt, err := time.Parse("2006-01-02 15:04:05", string(updatedAtStr))
			if err != nil {
				http.Error(w, "Failed to parse updated_at", http.StatusInternalServerError)
				log.Println("Failed to parse updated_at:", err)
				return
			}
			note.UpdatedAt = updatedAt
			notes = append(notes, note)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Failed to process rows", http.StatusInternalServerError)
			return
		}

		// HTML表示
		tmpl := template.Must(template.ParseFiles("./template/index.html"))
		if err := tmpl.Execute(w, notes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
