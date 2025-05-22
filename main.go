package main

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type PageData struct {
	Title   string
	Message string
}

func init() {
	// Set zerolog default timestamp format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Optional: pretty output in development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

// func initDB() (*gorm.DB, error) {
// 	dsn := "host=localhost user=postgres password=root dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta"
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Error().Err(err).Msg("Failed to connect to database")
// 		return nil, fmt.Errorf("failed to connect to database: %v", err)
// 	}

// 	log.Info().Msg("Database connected successfully")
// 	return db, nil
// }

func main() {
	log.Info().Msg("Initializing routes...")
	
	http.HandleFunc("/", logMiddleware(homeHandler))

	http.HandleFunc("/api/message", logMiddleware(apiHandler))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Info().Msg("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

// Middleware logging untuk semua handler
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
		Title:   "Welcome to Go Web App",
		Message: "Hello from Go Backend!",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
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

	log.Info().
		Str("response", response).
		Msg("API response sent successfully")
}
