package location

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type LocationClient struct {
	client *http.Client
	tracer trace.Tracer
}

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

func NewLocationClient() *LocationClient {
	return &LocationClient{
		client: &http.Client{},
		tracer: otel.Tracer("location-client"),
	}
}

func (c *LocationClient) GetLocation(ctx context.Context, cep string) (string, error) {
	ctx, span := c.tracer.Start(ctx, "getLocation")
	defer span.End()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get location")
	}

	var locationResponse ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&locationResponse); err != nil {
		return "", err
	}

	if locationResponse.Localidade == "" {
		return "", fmt.Errorf("can not find zipcode")
	}

	return locationResponse.Localidade, nil
}
