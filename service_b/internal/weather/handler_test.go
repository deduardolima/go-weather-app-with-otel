package weather

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockLocationService struct{}

func (mls *MockLocationService) GetLocation(ctx context.Context, cep string) (string, error) {
	if cep == "12345678" {
		return "São Paulo", nil
	}
	return "", fmt.Errorf("can not find zipcode")
}

type MockWeatherService struct {
	returnError bool
}

func (mws *MockWeatherService) GetWeather(ctx context.Context, city string) (float64, error) {
	if mws.returnError {
		return 0, fmt.Errorf("failed to get weather data")
	}
	if city == "São Paulo" {
		return 28.5, nil
	}
	return 0, fmt.Errorf("failed to get weather data")
}

func TestWeatherHandler_ValidCEP(t *testing.T) {
	locationService := &MockLocationService{}
	weatherService := &MockWeatherService{returnError: false}
	handler := NewWeatherHandler(locationService, weatherService)

	input := map[string]string{"cep": "12345678"}
	jsonInput, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/weather", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response WeatherResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "São Paulo", response.City)
	assert.Equal(t, 28.5, response.TempC)
	assert.Equal(t, 83.30000000000001, response.TempF)
	assert.Equal(t, 301.65, response.TempK)
}

func TestWeatherHandler_InvalidRequestBody(t *testing.T) {
	locationService := &MockLocationService{}
	weatherService := &MockWeatherService{}
	handler := NewWeatherHandler(locationService, weatherService)

	req := httptest.NewRequest("POST", "/weather", bytes.NewBuffer([]byte("invalid body")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "invalid request body\n", rr.Body.String())
}

func TestWeatherHandler_InvalidCEP(t *testing.T) {
	locationService := &MockLocationService{}
	weatherService := &MockWeatherService{}
	handler := NewWeatherHandler(locationService, weatherService)

	input := map[string]string{"cep": "123"}
	jsonInput, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/weather", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, "invalid zipcode\n", rr.Body.String())
}

func TestWeatherHandler_NotFoundCEP(t *testing.T) {
	locationService := &MockLocationService{}
	weatherService := &MockWeatherService{}
	handler := NewWeatherHandler(locationService, weatherService)

	input := map[string]string{"cep": "87654321"}
	jsonInput, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/weather", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "can not find zipcode\n", rr.Body.String())
}

func TestWeatherHandler_FailedWeatherData(t *testing.T) {
	locationService := &MockLocationService{}
	weatherService := &MockWeatherService{returnError: true}
	handler := NewWeatherHandler(locationService, weatherService)

	input := map[string]string{"cep": "12345678"}
	jsonInput, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/weather", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "failed to get weather data\n", rr.Body.String())
}
