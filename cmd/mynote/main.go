package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Note struct {
	ID            int
	Title         string
	Content       string
	Category      string
	ImportantFlag bool
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
		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		// データの取得
		rows, err := db.Query("SELECT * FROM note")
		if err != nil {
			http.Error(w, "Failed to execute query", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// 結果を表示
		var notes []Note
		for rows.Next() {
			var note Note
			if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.Category, &note.ImportantFlag); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			notes = append(notes, note)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Failed to process rows", http.StatusInternalServerError)
			return
		}

		// HTML表示
		tmpl := template.Must(template.ParseFiles("index.html"))
		if err := tmpl.Execute(w, notes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
