package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PageData struct {
	Title   string
	Message string
}

type Contact struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string
	Message   string
	CreatedAt time.Time
}

var db *gorm.DB

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	DB, err := initDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	db = DB
}

func initDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=root dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	log.Info().Msg("Database connected successfully")

	if err := db.AutoMigrate(&Contact{}); err != nil {
		log.Error().Err(err).Msg("AutoMigrate failed")
		return nil, err
	}
	log.Info().Msg("AutoMigrate executed successfully")

	return db, nil
}

func main() {
	log.Info().Msg("Initializing routes...")

	http.HandleFunc("/", logMiddleware(homeHandler))
	http.HandleFunc("/api/message", logMiddleware(apiHandler))
	http.HandleFunc("/contact", logMiddleware(contactHandler))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Info().Msg("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Msg("Incoming request")
		next(w, r)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Error().Err(err).Msg("Error loading template")
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:   "Welcome to My Portfolio",
		Message: "Hello from Go Backend!",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Error().Err(err).Msg("Error rendering template")
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	log.Info().Msg("Home page rendered successfully")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := `{"message": "This is an API response"}`
	w.Write([]byte(response))

	log.Info().Str("response", response).Msg("API response sent successfully")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Error().Err(err).Msg("Failed to parse form")
		http.Error(w, "Form parse error", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	message := r.FormValue("message")

	if name == "" || email == "" || message == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	contact := Contact{
		Name:      name,
		Email:     email,
		Message:   message,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&contact).Error; err != nil {
		log.Error().Err(err).Msg("Failed to save contact")
		http.Error(w, "Failed to save contact", http.StatusInternalServerError)
		return
	}

	log.Info().Str("email", email).Msg("Contact saved successfully")
	http.Redirect(w, r, "/?sent=true", http.StatusSeeOther)
}
