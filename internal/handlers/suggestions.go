package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"strconv"

	"groupie-tracker/internal/cache"
	"groupie-tracker/internal/models"
)

// HandleSuggestions handles the suggestions API endpoint
func HandleSuggestions(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    if query == "" {
        ErrorHandler(w, r, http.StatusBadRequest, "Missing search query")
        return
    }

    cachedData, err := cache.GetCachedData()
    if err != nil {
        ErrorHandler(w, r, http.StatusInternalServerError, "Failed to fetch data")
        return
    }

    suggestions := getSuggestions(query, cachedData)
    json.NewEncoder(w).Encode(suggestions)
}

func getSuggestions(query string, cachedData models.Datas) []models.Suggestion {
    var suggestions []models.Suggestion
    lowercaseQuery := strings.ToLower(query)
    
    // Use a map to prevent duplicate suggestions
    uniqueSuggestions := make(map[string]bool)

    // Create a map of artist ID to locations for quick lookup
    artistLocations := make(map[int][]string)
    for _, loc := range cachedData.LocationsData.Index {
        artistLocations[loc.ID] = loc.Locations
    }

    for _, artist := range cachedData.ArtistsData {
        // Artist name suggestions
        if strings.Contains(strings.ToLower(artist.Name), lowercaseQuery) {
            addUniqueSuggestion(&suggestions, uniqueSuggestions, artist.Name, "artist/band")
        }

        // Member suggestions
        for _, member := range artist.Members {
            if strings.Contains(strings.ToLower(member), lowercaseQuery) {
                addUniqueSuggestion(&suggestions, uniqueSuggestions, member, "member")
            }
        }

        // First album suggestions
        if strings.Contains(strings.ToLower(artist.FirstAlbum), lowercaseQuery) {
            addUniqueSuggestion(&suggestions, uniqueSuggestions, artist.FirstAlbum, "first album")
        }

        // Creation date suggestions
        if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
            addUniqueSuggestion(&suggestions, uniqueSuggestions, strconv.Itoa(artist.CreationDate), "created date")
        }

        // Location suggestions
        locations := artistLocations[artist.ID]
        for _, location := range locations {
            location = strings.TrimSpace(location)
            if strings.Contains(strings.ToLower(location), lowercaseQuery) {
                addUniqueSuggestion(&suggestions, uniqueSuggestions, location, "location")
            }
        }
    }

    return suggestions
}

// addUniqueSuggestion adds a suggestion if it hasn't been added before
func addUniqueSuggestion(suggestions *[]models.Suggestion, unique map[string]bool, text, sugType string) {
    key := text + "|" + sugType
    if !unique[key] {
        *suggestions = append(*suggestions, models.Suggestion{Text: text, Type: sugType})
        unique[key] = true
    }
}