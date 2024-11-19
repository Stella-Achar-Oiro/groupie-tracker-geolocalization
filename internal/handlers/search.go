package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "groupie-tracker/internal/cache"
    "groupie-tracker/internal/models"
)

// HandleSearch handles the search API endpoint
func HandleSearch(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    var filters models.FilterParams
    err := json.NewDecoder(r.Body).Decode(&filters)
    if err != nil {
        ErrorHandler(w, r, http.StatusBadRequest, "Invalid filter parameters")
        return
    }

    cachedData, err := cache.GetCachedData()
    if err != nil {
        ErrorHandler(w, r, http.StatusInternalServerError, "Failed to fetch data")
        return
    }

    results := searchArtists(query, cachedData, filters)
    json.NewEncoder(w).Encode(results)
}

func searchArtists(query string, cachedData models.Datas, filters models.FilterParams) models.SearchResult {
    var results models.SearchResult
    lowercaseQuery := strings.ToLower(query)

    // Try to parse query as a year
    queryYear, err := strconv.Atoi(query)
    isYearQuery := err == nil && len(query) == 4

    // Create a map of artist ID to locations for quick lookup
    artistLocations := make(map[int][]string)
    for _, loc := range cachedData.LocationsData.Index {
        artistLocations[loc.ID] = loc.Locations
    }

    for _, artist := range cachedData.ArtistsData {
        // Get locations for this artist
        locations := artistLocations[artist.ID]
        
        // Check if artist matches filters first
        if !matchesFilters(artist, locations, filters) {
            continue
        }

        // If query is empty, include all artists that match filters
        if query == "" {
            results.Artists = append(results.Artists, artist)
            continue
        }

        var matches bool
        
        // If it's a year query, ONLY check the creation date
        if isYearQuery {
            matches = (queryYear == artist.CreationDate)
        } else {
            // For non-year queries, check other fields
            matches = strings.Contains(strings.ToLower(artist.Name), lowercaseQuery) ||
                     containsAny(artist.Members, lowercaseQuery) ||
                     strings.Contains(strings.ToLower(artist.FirstAlbum), lowercaseQuery) ||
                     containsLocation(locations, lowercaseQuery)
        }

        if matches {
            results.Artists = append(results.Artists, artist)
        }
    }
    return results
}

// containsLocation checks if any location contains the search query
func containsLocation(locations []string, query string) bool {
    for _, location := range locations {
        if strings.Contains(strings.ToLower(strings.TrimSpace(location)), query) {
            return true
        }
    }
    return false
}

func containsAny(slice []string, substr string) bool {
    for _, s := range slice {
        if strings.Contains(strings.ToLower(s), substr) {
            return true
        }
    }
    return false
}

func matchesFilters(artist models.Artist, locations []string, filters models.FilterParams) bool {
    // Check creation year
    if artist.CreationDate < filters.CreationYearMin || artist.CreationDate > filters.CreationYearMax {
        return false
    }

    // Check first album year
    firstAlbumYear, _ := strconv.Atoi(strings.Split(artist.FirstAlbum, "-")[2])
    if firstAlbumYear < filters.FirstAlbumYearMin || firstAlbumYear > filters.FirstAlbumYearMax {
        return false
    }

    // Check number of members
    if len(filters.Members) > 0 {
        memberCount := len(artist.Members)
        if !contains(filters.Members, memberCount) {
            return false
        }
    }

    // Check locations
    if len(filters.Locations) > 0 {
        matched := false
        for _, loc := range locations {
            loc = strings.TrimSpace(loc)
            for _, filterLoc := range filters.Locations {
                if strings.Contains(strings.ToLower(loc), strings.ToLower(filterLoc)) {
                    matched = true
                    break
                }
            }
            if matched {
                break
            }
        }
        if !matched {
            return false
        }
    }

    return true
}

// contains checks if a slice contains a value
func contains(slice []int, val int) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}