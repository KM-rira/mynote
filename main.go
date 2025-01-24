package main

import (
	"os"

	"mynote/internal/handlers"
	"mynote/internal/model"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initLogger() {
	// app.logファイルを開く（存在しない場合は作成）
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("ファイルを開けませんでした: %v", err)
	}

	// logrusの出力先をファイルに設定
	log.SetOutput(file)

	// ログフォーマットの設定（オプション）
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// ログレベルの設定（必要に応じて変更）
	log.SetLevel(log.InfoLevel)
}

func main() {
	initLogger()

	// GORM を用いて SQLite データベースに接続（存在しない場合は新規作成）
	db, err := gorm.Open(sqlite.Open("mynote.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// モデルに基づく自動マイグレーションでテーブルを作成
	err = db.AutoMigrate(&model.Note{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// サンプルデータの挿入（存在しない場合にのみ追加）
	notes := []model.Note{
		{Title: "Sample Note 1", Contents: "This is the content of sample note 1", Category: "General", Important: true},
		{Title: "Sample Note 2", Contents: "This is the content of sample note 2", Category: "Work", Important: false},
		{Title: "Sample Note 3", Contents: "This is the content of sample note 3", Category: "Personal", Important: true},
		{Title: "サンプルデータ4です。", Contents: "こんにちは", Category: "テスト", Important: true},
	}

	for _, note := range notes {
		err = db.FirstOrCreate(&note, model.Note{Title: note.Title}).Error
		if err != nil {
			log.Fatalf("Failed to insert sample data: %v", err)
		}
	}

	// Handler インスタンスの生成
	h := handlers.NewHandler(db)
	// Echo インスタンスの作成
	e := echo.New()

	// 静的ファイルの提供設定
	e.Static("/css", "internal/templates/css")
	e.Static("/js", "internal/templates/js")

	// ルートの定義
	e.GET("/", func(c echo.Context) error {
		h.Index(c.Response().Writer, c.Request())
		return nil
	})
	// 他のルートも同様に設定
	e.GET("/register-form", func(c echo.Context) error {
		h.RegisterForm(c.Response().Writer, c.Request())
		return nil
	})
	e.Any("/register", func(c echo.Context) error {
		h.Register(c.Response().Writer, c.Request())
		return nil
	})
	e.GET("/select", func(c echo.Context) error {
		h.Select(c.Response().Writer, c.Request())
		return nil
	})
	e.Any("/update", func(c echo.Context) error {
		h.Update(c.Response().Writer, c.Request())
		return nil
	})
	e.Any("/delete", func(c echo.Context) error {
		h.Delete(c.Response().Writer, c.Request())
		return nil
	})

	log.Println("Server is running on port 9999")
	if err := e.Start(":9999"); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
