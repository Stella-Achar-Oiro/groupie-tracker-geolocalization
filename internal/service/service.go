package service

const (
    ArtistsAPI   = "https://groupietrackers.herokuapp.com/api/artists"
    LocationsAPI = "https://groupietrackers.herokuapp.com/api/locations"
    DatesAPI     = "https://groupietrackers.herokuapp.com/api/dates"
    RelationsAPI = "https://groupietrackers.herokuapp.com/api/relation"
)

const (
    MapboxAccessToken  = "pk.eyJ1Ijoic3RlbGxhYWNoYXJvaXJvIiwiYSI6ImNtMWhmZHNlODBlc3cybHF5OWh1MDI2dzMifQ.wk3v-v7IuiSiPwyq13qdHw"
    MapboxGeocodingAPI = "https://api.mapbox.com/geocoding/v5/mapbox.places"
)

func GetArtistsAPI() string {
    return ArtistsAPI
}

func GetLocationsAPI() string {
    return LocationsAPI
}

func GetDatesAPI() string {
    return DatesAPI
}

func GetRelationsAPI() string {
    return RelationsAPI
}

func GetMapboxAccessToken() string {
    return MapboxAccessToken
}

func GetMapboxGeocodingAPI() string {
    return MapboxGeocodingAPI
}
