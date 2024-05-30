package location

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
}

type LocationClient struct {
	tracer trace.Tracer
}

func NewLocationClient(tracer trace.Tracer) *LocationClient {
	return &LocationClient{tracer: tracer}
}

func (lc *LocationClient) GetLocation(ctx context.Context, cep string) (*ViaCEP, error) {
	_, span := lc.tracer.Start(ctx, "LocationClient_GetLocation")
	defer span.End()

	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch location for CEP: %s", cep)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var viaCEP ViaCEP
	err = json.Unmarshal(body, &viaCEP)
	if err != nil {
		return nil, err
	}

	if viaCEP.Localidade == "" {
		return nil, fmt.Errorf("can not find zipcode")
	}

	return &viaCEP, nil
}
