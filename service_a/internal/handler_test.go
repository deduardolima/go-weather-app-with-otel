package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRoundTripper struct {
	StatusCode int
	Response   string
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	response := &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(bytes.NewBufferString(m.Response)),
		Header:     make(http.Header),
	}
	return response, nil
}

func TestInputHandler_ValidCEP(t *testing.T) {
	os.Setenv("SERVICE_B_URL", "http://service-b:8081/weather")

	input := Input{CEP: "12345678"}
	jsonInput, _ := json.Marshal(input)

	mockResponse := `{"city":"Curitiba","temp_C":20,"temp_F":68,"temp_K":293.15}`
	mockTransport := &MockRoundTripper{StatusCode: http.StatusOK, Response: mockResponse}
	client := &http.Client{Transport: mockTransport}

	req := httptest.NewRequest("POST", "/input", bytes.NewBuffer(jsonInput))
	rr := httptest.NewRecorder()

	handler := NewInputHandler(client)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, mockResponse, rr.Body.String())
}

func TestInputHandler_InvalidCEP(t *testing.T) {
	input := Input{CEP: "12345"}
	jsonInput, _ := json.Marshal(input)

	req := httptest.NewRequest("POST", "/input", bytes.NewBuffer(jsonInput))
	rr := httptest.NewRecorder()

	handler := NewInputHandler(http.DefaultClient)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, "invalid zipcode\n", rr.Body.String())
}

func TestInputHandler_FailedToCreateRequest(t *testing.T) {
	os.Setenv("SERVICE_B_URL", "http://service-b:8081/weather")

	input := Input{CEP: "12345678"}
	jsonInput, _ := json.Marshal(input)

	req := httptest.NewRequest("POST", "/input", bytes.NewBuffer(jsonInput))
	rr := httptest.NewRecorder()

	mockTransport := &MockRoundTripper{
		StatusCode: http.StatusInternalServerError,
		Response:   "failed to create request",
	}
	client := &http.Client{Transport: mockTransport}

	handler := NewInputHandler(client)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "failed to create request\n", rr.Body.String())
}

func TestInputHandler_FailedToGetResponseFromServiceB(t *testing.T) {
	os.Setenv("SERVICE_B_URL", "http://service-b:8081/weather")

	input := Input{CEP: "12345678"}
	jsonInput, _ := json.Marshal(input)

	req := httptest.NewRequest("POST", "/input", bytes.NewBuffer(jsonInput))
	rr := httptest.NewRecorder()

	mockTransport := &MockRoundTripper{StatusCode: http.StatusInternalServerError, Response: "failed to get response from service B"}
	client := &http.Client{Transport: mockTransport}

	handler := NewInputHandler(client)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "failed to get response from service B\n", rr.Body.String())
}

func TestInputHandler_ServiceBErrorResponse(t *testing.T) {
	os.Setenv("SERVICE_B_URL", "http://service-b:8081/weather")

	input := Input{CEP: "12345678"}
	jsonInput, _ := json.Marshal(input)

	mockResponse := `{"error":"zipcode not found"}`
	mockTransport := &MockRoundTripper{StatusCode: http.StatusNotFound, Response: mockResponse}
	client := &http.Client{Transport: mockTransport}

	req := httptest.NewRequest("POST", "/input", bytes.NewBuffer(jsonInput))
	rr := httptest.NewRecorder()

	handler := NewInputHandler(client)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.JSONEq(t, mockResponse, rr.Body.String())
}
