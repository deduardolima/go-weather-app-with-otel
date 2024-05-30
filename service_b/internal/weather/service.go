package weather

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type LocationService interface {
	GetLocation(ctx context.Context, cep string) (string, error)
}

type WeatherService struct {
	client *WeatherClient
	tracer trace.Tracer
}

func NewWeatherService(tracer trace.Tracer) *WeatherService {
	return &WeatherService{
		client: NewWeatherClient(tracer),
		tracer: tracer,
	}
}

func (ws *WeatherService) GetWeather(ctx context.Context, city string) (float64, error) {
	ctx, span := ws.tracer.Start(ctx, "WeatherService_GetWeather")
	defer span.End()

	return ws.client.GetWeather(ctx, city)
}
