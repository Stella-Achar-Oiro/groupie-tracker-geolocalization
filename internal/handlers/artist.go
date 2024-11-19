package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"groupie-tracker/internal/cache"
	"groupie-tracker/internal/models"
)

func HandleArtist(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest, "Invalid artist ID")
		return
	}

	cachedData, err := cache.GetCachedData()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to fetch data")
		return
	}

	var artist models.Artist
	for _, a := range cachedData.ArtistsData {
		if a.ID == id {
			artist = a
			break
		}
	}

	if artist.ID == 0 {
		ErrorHandler(w, r, http.StatusNotFound, "Artist not found")
		return
	}

	details := models.ArtistDetail{
		Artist:    artist,
		Locations: getLocations(id, cachedData.LocationsData),
		Dates:     getDates(id, cachedData.DatesData),
		Relations: getRelations(id, cachedData.RelationsData),
	}

	json.NewEncoder(w).Encode(details)
}

func getLocations(id int, locationsData models.Location) []models.GeoLocation {
	for _, loc := range locationsData.Index {
		if loc.ID == id {
			var geoLocations []models.GeoLocation
			for _, location := range loc.Locations {
				geoLoc, err := geocode(location)
				if err != nil {
					log.Printf("Failed to geocode location: %v", err)
					continue
				}
				geoLocations = append(geoLocations, geoLoc)
			}
			return geoLocations
		}
	}
	return nil
}

func getDates(id int, datesData models.Date) []string {
	for _, date := range datesData.Index {
		if date.ID == id {
			return date.Dates
		}
	}
	return nil
}

func getRelations(id int, relationsData models.Relation) map[string][]string {
	for _, rel := range relationsData.Index {
		if rel.ID == id {
			return rel.DatesLocations
		}
	}
	return nil
}

func geocode(address string) (models.GeoLocation, error) {
	mapboxGeocodingAPI := models.GetMapboxGeocodingAPI()
	mapboxAccessToken := models.GetMapboxAccessToken()
	
	url := fmt.Sprintf("%s/%s.json?access_token=%s", mapboxGeocodingAPI, url.QueryEscape(address), mapboxAccessToken)

	resp, err := http.Get(url)
	if err != nil {
		return models.GeoLocation{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Features []struct {
			Center [2]float64 `json:"center"`
		} `json:"features"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.GeoLocation{}, err
	}

	if len(result.Features) == 0 {
		return models.GeoLocation{}, fmt.Errorf("no results found for address: %s", address)
	}

	return models.GeoLocation{
		Address: address,
		Lon:     result.Features[0].Center[0],
		Lat:     result.Features[0].Center[1],
	}, nil
}