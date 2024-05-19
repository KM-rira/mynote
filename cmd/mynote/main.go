package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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

	loc := url.QueryEscape("Asia/Tokyo")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s", dbUser, dbPassword, dbHost, dbPort, dbName, loc)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}

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
			var updatedAtStr string
			if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.Category, &note.ImportantFlag, &updatedAtStr); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				log.Println("Failed to scan row:", err)
				return
			}
			updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
			if err != nil {
				http.Error(w, "Failed to parse updated_at", http.StatusInternalServerError)
				log.Println("Failed to parse updated_at:", err)
				return
			}
			note.UpdatedAt = updatedAt.In(jst)
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

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("access /add")
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")
		important := r.FormValue("important") == "on"

		_, err := db.Exec("INSERT INTO note (title, content, category, important_flag) VALUES (?, ?, ?, ?)", title, content, category, important)
		if err != nil {
			http.Error(w, "Failed to insert note", http.StatusInternalServerError)
			log.Println("Failed to insert note:", err)
			return
		}

		http.Redirect(w, r, "/index", http.StatusSeeOther)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("access /delete")
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()
		ids := r.Form["select"]
		if len(ids) == 0 {
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		}

		query := fmt.Sprintf("DELETE FROM note WHERE id IN (%s)", strings.Join(ids, ","))
		_, err := db.Exec(query)
		if err != nil {
			http.Error(w, "Failed to delete notes", http.StatusInternalServerError)
			log.Println("Failed to delete notes:", err)
			return
		}

		http.Redirect(w, r, "/index", http.StatusSeeOther)
	})
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
