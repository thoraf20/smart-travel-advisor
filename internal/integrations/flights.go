package integrations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/thoraf20/smart-travel-advisor/internal/cache"
)

type FlightData struct {
	Airline struct {
		Name string `json:"name"`
	} `json:"airline"`
	Flight struct {
		Number string `json:"number"`
	} `json:"flight"`
	Departure struct {
		Airport string `json:"airport"`
		Time    string `json:"scheduled"`
	} `json:"departure"`
	Arrival struct {
		Airport string `json:"airport"`
		Time    string `json:"scheduled"`
	} `json:"arrival"`
}

type FlightResponse struct {
	Data []FlightData `json:"data"`
}

func GetFlightsArrivingInCity(cityName string) ([]FlightData, error) {
	apiKey := viper.GetString("FLIGHT_API_KEY")
	url := fmt.Sprintf("http://api.aviationstack.com/v1/flights?access_key=%s&arr_iata=%s&limit=3", apiKey, cityName)

	cacheKey := "flights:" + cityName

	cached, err := cache.CacheGet(cacheKey)
	if err == nil && cached != "" {
		var cachedResp FlightResponse
		if err := json.Unmarshal([]byte(cached), &cachedResp); 
		err == nil {
			return cachedResp.Data, nil
		}
	}

	var result FlightResponse
	resp, err := resty.New().R().
		SetResult(&result).
		Get(url)

	fmt.Println("Flight API response:", resp, err)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("flight API returned %s", resp.Status())
	}

	if data, err := json.Marshal(result);
	err == nil {
		cache.CacheSet(cacheKey, string(data), time.Hour)
	}

	return result.Data, nil
}
