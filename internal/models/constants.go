package models

var (
	MapboxAccessToken  string
	MapboxGeocodingAPI string
)

const (
	 // Artist ID limits
	 MinArtistID = 1
	 MaxArtistID = 52
)

// InitConstants initializes the package-level constants
func InitConstants(accessToken, geocodingAPI string) {
	MapboxAccessToken = accessToken
	MapboxGeocodingAPI = geocodingAPI
}

// GetMapboxAccessToken returns the Mapbox access token
func GetMapboxAccessToken() string {
	return MapboxAccessToken
}

// GetMapboxGeocodingAPI returns the Mapbox Geocoding API URL
func GetMapboxGeocodingAPI() string {
	return MapboxGeocodingAPI
}