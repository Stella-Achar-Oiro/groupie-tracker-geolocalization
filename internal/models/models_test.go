package models

import (
	"encoding/json"
	"testing"
	"time"
)

// TestArtistJSON tests the Artist struct JSON marshaling/unmarshaling
func TestArtistJSON(t *testing.T) {
	artistJSON := `{
		"id": 1,
		"image": "test.jpg",
		"name": "Test Artist",
		"members": ["Member 1", "Member 2"],
		"creationDate": 2020,
		"firstAlbum": "2020-01-01",
		"locations": "test-locations",
		"concertDates": "test-dates",
		"relations": "test-relations"
	}`

	var artist Artist
	if err := json.Unmarshal([]byte(artistJSON), &artist); err != nil {
		t.Fatalf("Failed to unmarshal Artist JSON: %v", err)
	}

	// Verify fields
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"ID", artist.ID, 1},
		{"Image", artist.Image, "test.jpg"},
		{"Name", artist.Name, "Test Artist"},
		{"CreationDate", artist.CreationDate, 2020},
		{"FirstAlbum", artist.FirstAlbum, "2020-01-01"},
		{"Locations", artist.Locations, "test-locations"},
		{"ConcertDates", artist.ConcertDates, "test-dates"},
		{"Relations", artist.Relations, "test-relations"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Artist.%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}

	// Check members array
	expectedMembers := []string{"Member 1", "Member 2"}
	if len(artist.Members) != len(expectedMembers) {
		t.Errorf("Artist.Members length = %d, want %d", len(artist.Members), len(expectedMembers))
	}
	for i, member := range artist.Members {
		if member != expectedMembers[i] {
			t.Errorf("Artist.Members[%d] = %s, want %s", i, member, expectedMembers[i])
		}
	}
}

// TestLocationJSON tests the Location struct JSON marshaling/unmarshaling
func TestLocationJSON(t *testing.T) {
	locationJSON := `{
		"index": [
			{
				"id": 1,
				"locations": ["New York, USA", "London, UK"],
				"dates": "test-dates"
			}
		]
	}`

	var location Location
	if err := json.Unmarshal([]byte(locationJSON), &location); err != nil {
		t.Fatalf("Failed to unmarshal Location JSON: %v", err)
	}

	if len(location.Index) != 1 {
		t.Fatalf("Location.Index length = %d, want 1", len(location.Index))
	}

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"ID", location.Index[0].ID, 1},
		{"Dates", location.Index[0].Dates, "test-dates"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Location.Index[0].%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}

	// Check locations array
	expectedLocations := []string{"New York, USA", "London, UK"}
	if len(location.Index[0].Locations) != len(expectedLocations) {
		t.Errorf("Location.Index[0].Locations length = %d, want %d",
			len(location.Index[0].Locations), len(expectedLocations))
	}
	for i, loc := range location.Index[0].Locations {
		if loc != expectedLocations[i] {
			t.Errorf("Location.Index[0].Locations[%d] = %s, want %s", i, loc, expectedLocations[i])
		}
	}
}

// TestDateJSON tests the Date struct JSON marshaling/unmarshaling
func TestDateJSON(t *testing.T) {
	dateJSON := `{
		"index": [
			{
				"id": 1,
				"dates": ["2020-01-01", "2020-02-01"]
			}
		]
	}`

	var date Date
	if err := json.Unmarshal([]byte(dateJSON), &date); err != nil {
		t.Fatalf("Failed to unmarshal Date JSON: %v", err)
	}

	if len(date.Index) != 1 {
		t.Fatalf("Date.Index length = %d, want 1", len(date.Index))
	}

	if date.Index[0].ID != 1 {
		t.Errorf("Date.Index[0].ID = %d, want 1", date.Index[0].ID)
	}

	expectedDates := []string{"2020-01-01", "2020-02-01"}
	if len(date.Index[0].Dates) != len(expectedDates) {
		t.Errorf("Date.Index[0].Dates length = %d, want %d",
			len(date.Index[0].Dates), len(expectedDates))
	}
	for i, d := range date.Index[0].Dates {
		if d != expectedDates[i] {
			t.Errorf("Date.Index[0].Dates[%d] = %s, want %s", i, d, expectedDates[i])
		}
	}
}

// TestRelationJSON tests the Relation struct JSON marshaling/unmarshaling
func TestRelationJSON(t *testing.T) {
	relationJSON := `{
		"index": [
			{
				"id": 1,
				"datesLocations": {
					"New York, USA": ["2020-01-01", "2020-02-01"],
					"London, UK": ["2020-03-01"]
				}
			}
		]
	}`

	var relation Relation
	if err := json.Unmarshal([]byte(relationJSON), &relation); err != nil {
		t.Fatalf("Failed to unmarshal Relation JSON: %v", err)
	}

	if len(relation.Index) != 1 {
		t.Fatalf("Relation.Index length = %d, want 1", len(relation.Index))
	}

	if relation.Index[0].ID != 1 {
		t.Errorf("Relation.Index[0].ID = %d, want 1", relation.Index[0].ID)
	}

	expectedDatesLocations := map[string][]string{
		"New York, USA": {"2020-01-01", "2020-02-01"},
		"London, UK":    {"2020-03-01"},
	}

	for location, dates := range expectedDatesLocations {
		gotDates, exists := relation.Index[0].DatesLocations[location]
		if !exists {
			t.Errorf("Location %s not found in DatesLocations", location)
			continue
		}
		if len(gotDates) != len(dates) {
			t.Errorf("DatesLocations[%s] length = %d, want %d",
				location, len(gotDates), len(dates))
			continue
		}
		for i, date := range dates {
			if gotDates[i] != date {
				t.Errorf("DatesLocations[%s][%d] = %s, want %s",
					location, i, gotDates[i], date)
			}
		}
	}
}

// TestSearchResultJSON tests the SearchResult struct JSON marshaling/unmarshaling
func TestSearchResultJSON(t *testing.T) {
	searchResultJSON := `{
		"artists": [
			{
				"id": 1,
				"name": "Test Artist",
				"members": ["Member 1"],
				"creationDate": 2020,
				"firstAlbum": "2020-01-01"
			}
		]
	}`

	var result SearchResult
	if err := json.Unmarshal([]byte(searchResultJSON), &result); err != nil {
		t.Fatalf("Failed to unmarshal SearchResult JSON: %v", err)
	}

	if len(result.Artists) != 1 {
		t.Fatalf("SearchResult.Artists length = %d, want 1", len(result.Artists))
	}

	artist := result.Artists[0]
	expected := Artist{
		ID:           1,
		Name:         "Test Artist",
		Members:      []string{"Member 1"},
		CreationDate: 2020,
		FirstAlbum:   "2020-01-01",
	}

	if artist.ID != expected.ID ||
		artist.Name != expected.Name ||
		artist.CreationDate != expected.CreationDate ||
		artist.FirstAlbum != expected.FirstAlbum {
		t.Errorf("SearchResult.Artists[0] = %+v, want %+v", artist, expected)
	}
}

// TestEventJSONMarshaling tests the Event struct JSON marshaling/unmarshaling
func TestEventJSONMarshaling(t *testing.T) {
	eventTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	event := Event{
		Location: "Test Location",
		Date:     eventTime,
	}

	// Test marshaling
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal Event: %v", err)
	}

	// Test unmarshaling
	var unmarshaledEvent Event
	if err := json.Unmarshal(data, &unmarshaledEvent); err != nil {
		t.Fatalf("Failed to unmarshal Event: %v", err)
	}

	if unmarshaledEvent.Location != event.Location {
		t.Errorf("Event Location = %s, want %s", unmarshaledEvent.Location, event.Location)
	}

	if !unmarshaledEvent.Date.Equal(event.Date) {
		t.Errorf("Event Date = %v, want %v", unmarshaledEvent.Date, event.Date)
	}
}

// TestFilterParamsJSON tests the FilterParams struct JSON marshaling/unmarshaling
func TestFilterParamsJSON(t *testing.T) {
	filterParamsJSON := `{
		"creationYearMin": 2000,
		"creationYearMax": 2020,
		"firstAlbumYearMin": 2000,
		"firstAlbumYearMax": 2020,
		"members": [1, 2, 3],
		"locations": ["New York", "London"]
	}`

	var params FilterParams
	if err := json.Unmarshal([]byte(filterParamsJSON), &params); err != nil {
		t.Fatalf("Failed to unmarshal FilterParams JSON: %v", err)
	}

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"CreationYearMin", params.CreationYearMin, 2000},
		{"CreationYearMax", params.CreationYearMax, 2020},
		{"FirstAlbumYearMin", params.FirstAlbumYearMin, 2000},
		{"FirstAlbumYearMax", params.FirstAlbumYearMax, 2020},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("FilterParams.%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}

	// Check members array
	expectedMembers := []int{1, 2, 3}
	if len(params.Members) != len(expectedMembers) {
		t.Errorf("FilterParams.Members length = %d, want %d",
			len(params.Members), len(expectedMembers))
	}
	for i, member := range params.Members {
		if member != expectedMembers[i] {
			t.Errorf("FilterParams.Members[%d] = %d, want %d", i, member, expectedMembers[i])
		}
	}

	// Check locations array
	expectedLocations := []string{"New York", "London"}
	if len(params.Locations) != len(expectedLocations) {
		t.Errorf("FilterParams.Locations length = %d, want %d",
			len(params.Locations), len(expectedLocations))
	}
	for i, location := range params.Locations {
		if location != expectedLocations[i] {
			t.Errorf("FilterParams.Locations[%d] = %s, want %s", i, location, expectedLocations[i])
		}
	}
}