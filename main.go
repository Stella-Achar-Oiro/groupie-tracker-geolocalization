package main

import (
    "html/template"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
    "groupie-tracker/internal/cache"
    "groupie-tracker/internal/handlers"
    "groupie-tracker/internal/models"
    "groupie-tracker/internal/service"
)

var (
    //indexTpl        *template.Template
    artistDetailsTpl *template.Template
    logger          *log.Logger
)

const cacheDuration = 1 * time.Hour

func main() {
    // Initialize logger
    logger = log.New(os.Stdout, "GROUPIE-TRACKER: ", log.Ldate|log.Ltime|log.Lshortfile)

    // Initialize models package with required constants
    models.InitConstants(service.GetMapboxAccessToken(), service.GetMapboxGeocodingAPI())
    logger.Println("Models initialized with Mapbox constants")

    // Initialize cache
    cache.Init(cacheDuration)
    logger.Println("Cache initialized with duration:", cacheDuration)

    // Initial data fetch
    if err := cache.RefreshCache(); err != nil {
        logger.Fatalf("Failed to fetch initial data: %v", err)
    }
    logger.Println("Initial data fetched successfully")

    // Serve index.html from the templates folder
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/static/") && r.URL.Path != "/" {
			handlers.ErrorHandler(w, r, http.StatusNotFound, "Page Not Found")
			return
		}
		// Check if index.html exists before proceeding
	tmplPath := "templates/index.html"
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		handlers.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to fetch data")
		return
	}
		http.ServeFile(w, r, "./templates/index.html")
	})

    // Parse HTML template
    var err error
    artistDetailsTpl, err = template.ParseFiles("templates/artist-details.html")
    if err != nil {
        logger.Fatalf("Failed to parse artist details template: %v", err)
    }

    // Set up routes
    http.HandleFunc("/artist/", enableCORS(handlers.HandleArtistDetails(artistDetailsTpl)))
    http.HandleFunc("/api/search", enableCORS(handlers.HandleSearch))
    http.HandleFunc("/api/artist/", enableCORS(handlers.HandleArtist))
    http.HandleFunc("/api/suggestions", enableCORS(handlers.HandleSuggestions))

    http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
        tmpl, err := template.ParseFiles("templates/about.html")
        if err != nil {
            handlers.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to load about page")
            return
        }
        tmpl.Execute(w, nil)
    })

    http.HandleFunc("/work-with-us", func(w http.ResponseWriter, r *http.Request) {
        tmpl, err := template.ParseFiles("templates/work-with-us.html")
        if err != nil {
            handlers.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to load work with us page")
            return
        }
        tmpl.Execute(w, nil)
    })

    // Serve static files
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    
    logger.Println("Routes and static file server set up")

    // Start server
    port := ":8080"
    logger.Printf("Server starting on %s", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        logger.Fatalf("Server failed to start: %v", err)
    }
}

func enableCORS(handler interface{}) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        // Handle different types of handlers
        switch h := handler.(type) {
        case http.HandlerFunc:
            h(w, r)
        case func(http.ResponseWriter, *http.Request):
            h(w, r)
        }
    }
}