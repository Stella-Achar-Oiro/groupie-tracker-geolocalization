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

    // Use the constants from models package
    if id < models.MinArtistID || id > models.MaxArtistID {
        ErrorHandler(w, r, http.StatusBadRequest, 
            fmt.Sprintf("Artist ID must be between %d and %d", models.MinArtistID, models.MaxArtistID))
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

    locations, err := getLocations(id, cachedData.LocationsData)
    if err != nil {
        log.Printf("Error getting locations: %v", err)
        ErrorHandler(w, r, http.StatusInternalServerError, "Failed to fetch locations")
        return
    }

    details := models.ArtistDetail{
        Artist:    artist,
        Locations: locations.Locations,
        Dates:     getDates(id, cachedData.DatesData),
        Relations: getRelations(id, cachedData.RelationsData),
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(details); err != nil {
        log.Printf("Error encoding response: %v", err)
        ErrorHandler(w, r, http.StatusInternalServerError, "Failed to encode response")
        return
    }
}

func getLocations(id int, locationsData models.Location) (*models.LocationPath, error) {
    var path models.LocationPath
    
    for _, loc := range locationsData.Index {
        if loc.ID == id {
            // Get coordinates for each location
            for i, location := range loc.Locations {
                geoLoc, err := geocode(location)
                if err != nil {
                    log.Printf("Failed to geocode location: %v", err)
                    continue
                }
                
                path.Locations = append(path.Locations, geoLoc)
                
                // Create path segments connecting consecutive locations
                if i > 0 {
                    prevLoc := path.Locations[i-1]
                    segment := models.PathSegment{
                        StartLat: prevLoc.Lat,
                        StartLon: prevLoc.Lon,
                        EndLat:   geoLoc.Lat,
                        EndLon:   geoLoc.Lon,
                    }
                    path.Path = append(path.Path, segment)
                }
            }
            return &path, nil
        }
    }
    return nil, fmt.Errorf("no locations found for artist ID: %d", id)
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
        return models.GeoLocation{}, fmt.Errorf("geocoding request failed: %v", err)
    }
    defer resp.Body.Close()

    var result struct {
        Features []struct {
            Center [2]float64 `json:"center"`
        } `json:"features"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return models.GeoLocation{}, fmt.Errorf("failed to decode geocoding response: %v", err)
    }

    if len(result.Features) == 0 {
        return models.GeoLocation{}, fmt.Errorf("no coordinates found for address: %s", address)
    }

    geoLoc := models.GeoLocation{
        Address: address,
        Lon:     result.Features[0].Center[0],
        Lat:     result.Features[0].Center[1],
    }

    log.Printf("Successfully geocoded %s to coordinates: [%f, %f]", 
        address, geoLoc.Lat, geoLoc.Lon)

    return geoLoc, nil
}