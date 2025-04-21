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

type WeatherResponse struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
}

func GetWeather(cityName string) (*WeatherResponse, error) {	
	apiKey := viper.GetString("WEATHER_API_KEY")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", cityName, apiKey)

	cacheKey := "weather:" + cityName

	cached, err := cache.CacheGet(cacheKey)
	if err == nil && cached != "" {
		var cachedResp WeatherResponse
		if err := json.Unmarshal([]byte(cached), &cachedResp); 
		err == nil {
			return &cachedResp, nil
		}
	}

	client := resty.New()
	var result WeatherResponse
	resp, err := client.R().
		SetResult(&result).
		Get(url)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("weather API returned %s", resp.Status())
	}

	if data, err := json.Marshal(result);
	err == nil {
		cache.CacheSet(cacheKey, string(data), time.Hour)
	}

	return &result, nil
}
