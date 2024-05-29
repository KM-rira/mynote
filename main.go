package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Note struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Contents  string `json:"contents"`
	Category  string `json:"category"`
	Important bool   `json:"important"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (Note) TableName() string {
	return "note"
}

func main() {
	// Load environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Data Source Name with charset and parseTime
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %s", err.Error())
	}
	gormDb, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err.Error())
	}

	// Template parsing
	tmplIndex := template.Must(template.ParseFiles("index.html"))
	tmplUpdate := template.Must(template.ParseFiles("update.html"))

	fmt.Println("Successfully connected to the database!")

	// HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/ Received request")

		var notes []Note
		result := gormDb.Order("updated_at DESC").Find(&notes)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInsufficientStorage)
			return
		}

		if result.RowsAffected == 0 {
			fmt.Fprintln(w, "No note record")
			return
		}

		// Render the template with notes
		log.Printf("Rendering template with notes: %+v", notes)
		err = tmplIndex.Execute(w, notes)
		if err != nil {
			log.Printf("Error rendering template: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/select", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/select Received request")
		id := r.URL.Query().Get("id")

		var note Note

		result := gormDb.First(&note, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		tmplUpdate.Execute(w, note)
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/update Received request")
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var note Note
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		note.ID, err = strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		note.Title = r.FormValue("title")
		note.Contents = r.FormValue("contents")
		note.Category = r.FormValue("category")
		note.Important = r.FormValue("important") == "on"

		_, err = db.Exec("UPDATE note SET title = ?, contents = ?, category = ?, important = ?, updated_at = NOW() WHERE id = ?", note.Title, note.Contents, note.Category, note.Important, note.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/register-form", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/register-form Received request")
		tmplRegister := template.Must(template.ParseFiles("register.html"))
		tmplRegister.Execute(w, nil)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/register Received request")
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var note Note
		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO note (title, contents, category, important, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())",
			note.Title, note.Contents, note.Category, note.Important)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/delete Received request")
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var requestData struct {
			ID int `json:"id"`
		}

		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		_, err = db.Exec("DELETE FROM note WHERE id = ?", requestData.ID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
