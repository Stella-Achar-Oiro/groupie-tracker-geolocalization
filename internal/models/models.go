package models

import "time"

const (
	ArtistsAPI   = "https://groupietrackers.herokuapp.com/api/artists"
	LocationsAPI = "https://groupietrackers.herokuapp.com/api/locations"
	DatesAPI     = "https://groupietrackers.herokuapp.com/api/dates"
	RelationsAPI = "https://groupietrackers.herokuapp.com/api/relation"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Location struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	} `json:"index"`
}

type Date struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

type Relation struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type Datas struct {
	ArtistsData   []Artist `json:"artists"`
	LocationsData Location `json:"locations"`
	DatesData     Date     `json:"dates"`
	RelationsData Relation `json:"relations"`
}

type SearchResult struct {
	Artists []Artist `json:"artists"`
}

type Suggestion struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type GeoLocation struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type Event struct {
	Location string    `json:"location"`
	Date     time.Time `json:"date"`
}

type ArtistDetail struct {
	Artist    Artist              `json:"artist"`
	Locations []GeoLocation       `json:"locations"`
	Dates     []string            `json:"dates"`
	Relations map[string][]string `json:"relations"`
	Events    []Event             `json:"events"`
}

type FilterParams struct {
	CreationYearMin   int      `json:"creationYearMin"`
	CreationYearMax   int      `json:"creationYearMax"`
	FirstAlbumYearMin int      `json:"firstAlbumYearMin"`
	FirstAlbumYearMax int      `json:"firstAlbumYearMax"`
	Members           []int    `json:"members"`
	Locations         []string `json:"locations"`
}
