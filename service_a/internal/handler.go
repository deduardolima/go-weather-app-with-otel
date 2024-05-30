package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

func NewCEPHandler(tracer trace.Tracer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CEPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.CEP) != 8 {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		ctx, span := tracer.Start(r.Context(), "ServiceA_HandleCEP")
		defer span.End()

		jsonData, _ := json.Marshal(req)

		client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
		httpReq, err := http.NewRequestWithContext(ctx, "POST", "http://service-b:8081/weather", bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "failed to create request for service B", http.StatusInternalServerError)
			return
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(httpReq)
		if err != nil {
			http.Error(w, "failed to communicate with service B", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			http.Error(w, string(body), resp.StatusCode)
			return
		}

		body, _ := io.ReadAll(resp.Body)
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}
