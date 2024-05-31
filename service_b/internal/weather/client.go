package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type WeatherClient struct {
	client *http.Client
	tracer trace.Tracer
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func NewWeatherClient() *WeatherClient {
	return &WeatherClient{
		client: &http.Client{},
		tracer: otel.Tracer("weather-client"),
	}
}

func (wc *WeatherClient) GetWeather(ctx context.Context, location string) (float64, error) {
	ctx, span := wc.tracer.Start(ctx, "GetWeather")
	defer span.End()

	apiKey := viper.GetString("WEATHER_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("WEATHER_API_KEY not set")
	}

	encodedLocation := url.QueryEscape(strings.ReplaceAll(location, " ", ""))
	apiURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, encodedLocation)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return 0, err
	}

	resp, err := wc.client.Do(req)
	if err != nil {
		log.Printf("Error performing request: %v\n", err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Response body: %s\n", string(body))
		return 0, fmt.Errorf("failed to get weather data")
	}

	var result WeatherAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.Current.TempC, nil
}
