package weather

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type WeatherHandler struct {
	locationService LocationService
	weatherService  WeatherService
	tracer          trace.Tracer
}

func NewWeatherHandler(locationService LocationService, weatherService WeatherService) *WeatherHandler {
	tracer := locationService.(interface{ Tracer() trace.Tracer }).Tracer()
	return &WeatherHandler{
		locationService: locationService,
		weatherService:  weatherService,
		tracer:          tracer,
	}
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type CEPRequest struct {
	CEP string `json:"cep"`
}

func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "WeatherHandler_ServeHTTP")
	defer span.End()

	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.CEP) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := h.locationService.GetLocation(ctx, req.CEP)
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
