package location

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type LocationService struct {
	client *LocationClient
	tracer trace.Tracer
}

func NewLocationService(tracer trace.Tracer) *LocationService {
	return &LocationService{
		client: NewLocationClient(tracer),
		tracer: tracer,
	}
}

func (ls *LocationService) GetLocation(ctx context.Context, cep string) (string, error) {
	ctx, span := ls.tracer.Start(ctx, "LocationService_GetLocation")
	defer span.End()

	viaCEP, err := ls.client.GetLocation(ctx, cep)
	if err != nil {
		return "", err
	}
	return viaCEP.Localidade, nil
}
