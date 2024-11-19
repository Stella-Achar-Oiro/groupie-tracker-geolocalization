package handlers

import (
    "html/template"
    "net/http"
    "strings"
)

func HandleArtistDetails(tpl *template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path
        artistID := strings.TrimPrefix(path, "/artist/")

        if artistID == "" {
            ErrorHandler(w, r, http.StatusBadRequest, "Missing artist ID")
            return
        }

        // Create template data with artist ID
        data := struct {
            ArtistID string
        }{
            ArtistID: artistID,
        }

        // Render the template
        err := tpl.Execute(w, data)
        if err != nil {
            ErrorHandler(w, r, http.StatusInternalServerError, "Failed to render template")
            return
        }
    }
}