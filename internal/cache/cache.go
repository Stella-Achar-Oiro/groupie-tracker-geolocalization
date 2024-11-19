package cache

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"groupie-tracker/internal/models"
	"groupie-tracker/internal/service"
)

var (
	cache    Cache
	duration time.Duration
)

type Cache struct {
	data      models.Datas
	expiresAt time.Time
	mutex     sync.RWMutex
}

func Init(cacheDuration time.Duration) {
	cache = Cache{}
	duration = cacheDuration
}

func RefreshCache() error {
	var newData models.Datas
	err := fetchAllData(&newData)
	if err != nil {
		return err
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.data = newData
	cache.expiresAt = time.Now().Add(duration)

	return nil
}

func GetCachedData() (models.Datas, error) {
	cache.mutex.RLock()
	if time.Now().Before(cache.expiresAt) {
		defer cache.mutex.RUnlock()
		return cache.data, nil
	}
	cache.mutex.RUnlock()

	if err := RefreshCache(); err != nil {
		return models.Datas{}, err
	}

	return cache.data, nil
}

func fetchAllData(data *models.Datas) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 4)

	wg.Add(4)
	go fetchData(service.GetArtistsAPI(), &data.ArtistsData, &wg, errChan)
	go fetchData(service.GetLocationsAPI(), &data.LocationsData, &wg, errChan)
	go fetchData(service.GetDatesAPI(), &data.DatesData, &wg, errChan)
	go fetchData(service.GetRelationsAPI(), &data.RelationsData, &wg, errChan)

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchData(url string, target interface{}, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		errChan <- fmt.Errorf("failed to fetch data from %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		errChan <- fmt.Errorf("failed to decode data from %s: %v", url, err)
	}
}
