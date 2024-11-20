package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

// Структура для данных о погоде от WeatherAPI.com
type WeatherData struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		Tz_id          string  `json:"tz_id"`
		LocalTimeEpoch int64   `json:"localtime_epoch"`
		LocalTime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int64   `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		Temp_c           float64 `json:"temp_c"`
		Temp_f           float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		Wind_mph    float64 `json:"wind_mph"`
		Wind_kph    float64 `json:"wind_kph"`
		Wind_degree int     `json:"wind_degree"`
		Wind_dir    string  `json:"wind_dir"`
		Pressure_mb float64 `json:"pressure_mb"`
		Pressure_in float64 `json:"pressure_in"`
		Precip_mm   float64 `json:"precip_mm"`
		Precip_in   float64 `json:"precip_in"`
		Humidity    int     `json:"humidity"`
		Cloud       int     `json:"cloud"`
		Feelslike_c float64 `json:"feelslike_c"`
		Feelslike_f float64 `json:"feelslike_f"`
		Windchill_c float64 `json:"windchill_c"`
		Windchill_f float64 `json:"windchill_f"`
		Heatindex_c float64 `json:"heatindex_c"`
		Heatindex_f float64 `json:"heatindex_f"`
		Dewpoint_c  float64 `json:"dewpoint_c"`
		Dewpoint_f  float64 `json:"dewpoint_f"`
		Vis_km      float64 `json:"vis_km"`
		Vis_miles   float64 `json:"vis_miles"`
		Uv          float64 `json:"uv"`
		Gust_mph    float64 `json:"gust_mph"`
		Gust_kph    float64 `json:"gust_kph"`
	} `json:"current"`
}

func getWeather(city string) (*WeatherData, error) {
	apiKey := "bc57f41501c241dc877190147242011"
	apiURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, url.QueryEscape(city))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неудачный запрос: %s", resp.Status)
	}

	var weatherData WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	return &weatherData, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.FormValue("city")
	if city == "" {
		city = "Samara" // Город по умолчанию
	}

	weather, err := getWeather(city)
	if err != nil {
		http.Error(w, "Ошибка получения погоды: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Ошибка обработки шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, weather)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", weatherHandler)
	fmt.Println("Сервер запущен на http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
