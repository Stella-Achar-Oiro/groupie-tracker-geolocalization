package service

import "testing"

func TestGetArtistsAPI(t *testing.T) {
	expected := "https://groupietrackers.herokuapp.com/api/artists"
	if got := GetArtistsAPI(); got != expected {
		t.Errorf("GetArtistsAPI() = %v, want %v", got, expected)
	}
}

func TestGetLocationsAPI(t *testing.T) {
	expected := "https://groupietrackers.herokuapp.com/api/locations"
	if got := GetLocationsAPI(); got != expected {
		t.Errorf("GetLocationsAPI() = %v, want %v", got, expected)
	}
}

func TestGetDatesAPI(t *testing.T) {
	expected := "https://groupietrackers.herokuapp.com/api/dates"
	if got := GetDatesAPI(); got != expected {
		t.Errorf("GetDatesAPI() = %v, want %v", got, expected)
	}
}

func TestGetRelationsAPI(t *testing.T) {
	expected := "https://groupietrackers.herokuapp.com/api/relation"
	if got := GetRelationsAPI(); got != expected {
		t.Errorf("GetRelationsAPI() = %v, want %v", got, expected)
	}
}

func TestGetMapboxAccessToken(t *testing.T) {
	expected := "pk.eyJ1Ijoic3RlbGxhYWNoYXJvaXJvIiwiYSI6ImNtMWhmZHNlODBlc3cybHF5OWh1MDI2dzMifQ.wk3v-v7IuiSiPwyq13qdHw"
	if got := GetMapboxAccessToken(); got != expected {
		t.Errorf("GetMapboxAccessToken() = %v, want %v", got, expected)
	}
}

func TestGetMapboxGeocodingAPI(t *testing.T) {
	expected := "https://api.mapbox.com/geocoding/v5/mapbox.places"
	if got := GetMapboxGeocodingAPI(); got != expected {
		t.Errorf("GetMapboxGeocodingAPI() = %v, want %v", got, expected)
	}
}