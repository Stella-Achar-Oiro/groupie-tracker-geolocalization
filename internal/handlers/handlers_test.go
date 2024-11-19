package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"groupie-tracker/internal/models"
)

// TestHandleArtistDetails tests the artist details handler
func TestHandleArtistDetails(t *testing.T) {
	tmpl := template.Must(template.New("artist-details").Parse(`{{.ArtistID}}`))

	tests := []struct {
		name           string
		urlPath        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid artist ID",
			urlPath:        "/artist/123",
			expectedStatus: http.StatusOK,
			expectedBody:   "123",
		},
		{
			name:           "Missing parts in path",
			urlPath:        "/artist",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid artist ID\n",
		},
		{
			name:           "Empty string ID",
			urlPath:        "/artist/",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.urlPath, nil)
			w := httptest.NewRecorder()

			handler := HandleArtistDetails(tmpl)
			handler(w, req)

			if got := w.Code; got != tt.expectedStatus {
				t.Errorf("HandleArtistDetails() status code = %v, want %v", got, tt.expectedStatus)
			}

			gotBody := w.Body.String()
			if gotBody != tt.expectedBody {
				t.Errorf("HandleArtistDetails() body = %q, want %q", gotBody, tt.expectedBody)
			}
		})
	}
}

// TestHandleArtist tests the artist API handler
func TestHandleArtist(t *testing.T) {
	tests := []struct {
		name           string
		urlPath        string
		expectedStatus int
	}{
		{
			name:           "Invalid artist ID format",
			urlPath:        "/api/artist/abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-existent artist ID",
			urlPath:        "/api/artist/999999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.urlPath, nil)
			w := httptest.NewRecorder()

			HandleArtist(w, req)

			if got := w.Code; got != tt.expectedStatus {
				t.Errorf("HandleArtist() status code = %v, want %v", got, tt.expectedStatus)
			}
		})
	}
}

// TestHandleSearch tests the search functionality
func TestHandleSearch(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		filters        models.FilterParams
		expectedStatus int
	}{
		{
			name:  "Valid search with filters",
			query: "test",
			filters: models.FilterParams{
				CreationYearMin:   2000,
				CreationYearMax:   2023,
				FirstAlbumYearMin: 2000,
				FirstAlbumYearMax: 2023,
				Members:           []int{1, 2, 3},
				Locations:         []string{"New York"},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "Empty query with filters",
			query: "",
			filters: models.FilterParams{
				CreationYearMin: 2000,
				CreationYearMax: 2023,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "Invalid filter format",
			query: "test",
			filters: models.FilterParams{
				CreationYearMin: 2025,
				CreationYearMax: 2000, // Invalid range
			},
			expectedStatus: http.StatusOK, // Still returns OK with empty results
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterJSON, err := json.Marshal(tt.filters)
			if err != nil {
				t.Fatalf("Failed to marshal filters: %v", err)
			}

			req := httptest.NewRequest("POST", "/api/search?q="+tt.query, bytes.NewBuffer(filterJSON))
			w := httptest.NewRecorder()

			HandleSearch(w, req)

			if got := w.Code; got != tt.expectedStatus {
				t.Errorf("HandleSearch() status code = %v, want %v", got, tt.expectedStatus)
			}

			if w.Code == http.StatusOK {
				var result models.SearchResult
				err := json.NewDecoder(w.Body).Decode(&result)
				if err != nil {
					t.Errorf("Failed to decode response body: %v", err)
				}
			}
		})
	}
}

// TestHandleSuggestions tests the suggestions functionality
func TestHandleSuggestions(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		expectedStatus int
	}{
		{
			name:           "Empty query",
			query:          "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Valid query",
			query:          "test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Special characters query",
			query:          "test@123",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/suggestions?q="+tt.query, nil)
			w := httptest.NewRecorder()

			HandleSuggestions(w, req)

			if got := w.Code; got != tt.expectedStatus {
				t.Errorf("HandleSuggestions() status code = %v, want %v", got, tt.expectedStatus)
			}

			if w.Code == http.StatusOK {
				var suggestions []models.Suggestion
				err := json.NewDecoder(w.Body).Decode(&suggestions)
				if err != nil {
					t.Errorf("Failed to decode response body: %v", err)
				}
			}
		})
	}
}

