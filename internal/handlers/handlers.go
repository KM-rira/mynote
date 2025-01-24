package handlers

import (
	"encoding/json"
	"html/template"
	"mynote/internal/model"
	"net/http"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}
type NoteResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Contents  string `json:"contents"`
	Category  string `json:"category"`
	Important bool   `json:"important"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// テンプレートの準備
	tmpl := template.Must(template.ParseFiles(filepath.Join("internal", "templates", "index.html")))

	// データベースからノートを取得
	var notes []model.Note
	result := h.db.Order("updated_at DESC").Find(&notes)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// ノートデータをフォーマット済みの構造体に変換
	var formattedNotes []NoteResponse
	for _, note := range notes {
		formattedNotes = append(formattedNotes, NoteResponse{
			ID:        note.ID,
			Title:     note.Title,
			Contents:  note.Contents,
			Category:  note.Category,
			Important: note.Important,
			CreatedAt: note.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: note.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// スライスをそのままテンプレートに渡す
	if err := tmpl.Execute(w, formattedNotes); err != nil {
		log.Errorf("テンプレートレンダリングエラー: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterForm(w http.ResponseWriter, r *http.Request) {
	// RegisterFormハンドラの実装
	log.Println("/register-form Received request")
	tmplRegister := template.Must(template.ParseFiles(filepath.Join("internal", "templates", "register.html")))
	tmplRegister.Execute(w, nil)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// Registerハンドラの実装
	log.Println("/register Received request")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var note model.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := h.db.Create(&note)
	if result.Error != nil {
		log.Printf("note record create error: %v\n", result.Error)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Select(w http.ResponseWriter, r *http.Request) {
	// Selectハンドラの実装
	log.Println("/select Received request")
	tmplUpdate := template.Must(template.ParseFiles(filepath.Join("internal", "templates", "update.html")))
	id := r.URL.Query().Get("id")

	var note model.Note

	result := h.db.First(&note, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	tmplUpdate.Execute(w, note)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Updateハンドラの実装
	log.Println("/update Received request")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var note model.Note
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

	result := h.db.Model(&note).Updates(model.Note{Title: note.Title, Contents: note.Contents, Category: note.Category, Important: note.Important})
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// Deleteハンドラの実装
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

	result := h.db.Delete(&model.Note{}, requestData.ID)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
