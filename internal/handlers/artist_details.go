package handlers

import (
	"fmt"
	"groupie-tracker/internal/models"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func HandleArtistDetails(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		artistIDStr := strings.TrimPrefix(path, "/artist/")

		if artistIDStr == "" {
			ErrorHandler(w, r, http.StatusBadRequest, "Missing artist ID")
			return
		}

		artistID, err := strconv.Atoi(artistIDStr)
		if err != nil {
			ErrorHandler(w, r, http.StatusBadRequest, "Invalid artist ID format")
			return
		}

		// Use the constants from models package
		if artistID < models.MinArtistID || artistID > models.MaxArtistID {
			ErrorHandler(w, r, http.StatusBadRequest,
				fmt.Sprintf("Artist ID must be between %d and %d", models.MinArtistID, models.MaxArtistID))
			return
		}

		// Create template data with artist ID
		data := struct {
			ArtistID string
		}{
			ArtistID: strconv.Itoa(artistID),
		}

		// Render the template
		err = tpl.Execute(w, data)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Failed to render template")
			return
		}
	}
}
