package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

type LocationService interface {
	GetLocation(ctx context.Context, cep string) (string, error)
}

type WeatherService interface {
	GetWeather(ctx context.Context, city string) (float64, error)
}

type WeatherHandler struct {
	locationService LocationService
	weatherService  WeatherService
}

func NewWeatherHandler(locationService LocationService, weatherService WeatherService) *WeatherHandler {
	return &WeatherHandler{
		locationService: locationService,
		weatherService:  weatherService,
	}
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("weather-handler")
	ctx, span := tr.Start(r.Context(), "ServeHTTP")
	defer span.End()

	fmt.Println("Mensagem recebida do servico A")

	var input struct {
		CEP string `json:"cep"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received input: %+v\n", input)

	if len(input.CEP) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	city, err := h.locationService.GetLocation(ctx, input.CEP)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	tempC, err := h.weatherService.GetWeather(ctx, city)
	if err != nil {
		http.Error(w, "failed to get weather data", http.StatusInternalServerError)
		return
	}

	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15

	response := WeatherResponse{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
