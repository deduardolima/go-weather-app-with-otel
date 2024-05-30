package internal_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/deduardolima/go-weather-with-otel/internal"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

type MockTracerProvider struct{}

type MockTracer struct{}

func (t *MockTracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {

	return ctx, nil
}

func (t *MockTracer) ForceFlush(ctx context.Context) error {

	return nil
}

func TestCEPHandler_ValidCEP(t *testing.T) {
	tracerProvider := &MockTracerProvider{}
	handler := internal.NewCEPHandler(tracerProvider)

	reqBody := `{"cep": "12345678"}`
	req := httptest.NewRequest("POST", "/cep", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponseBody := `{"temp_C": 20, "temp_F": 68, "temp_K": 293.15}`
	assert.Equal(t, expectedResponseBody, w.Body.String())
}

func TestCEPHandler_InvalidCEP(t *testing.T) {
	tracer := &MockTracerProvider{}
	handler := internal.NewCEPHandler(tracer)

	reqBody := `{"cep": "1234"}`
	req := httptest.NewRequest("POST", "/cep", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	expectedResponseBody := "invalid zipcode\n"
	assert.Equal(t, expectedResponseBody, w.Body.String())
}

func TestCEPHandler_ServiceBError(t *testing.T) {
	tracer := &MockTracerProvider{}
	handler := internal.NewCEPHandler(tracer)

	reqBody := `{"cep": "12345678"}`
	req := httptest.NewRequest("POST", "/cep", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	transport := &MockTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("service B error")
		},
	}
	client := &http.Client{Transport: transport}

	oldClient := http.DefaultClient
	defer func() { http.DefaultClient = oldClient }()
	http.DefaultClient = client

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedResponseBody := "failed to communicate with service B\n"
	assert.Equal(t, expectedResponseBody, w.Body.String())
}
