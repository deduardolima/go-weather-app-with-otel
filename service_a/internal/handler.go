package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
)

type Input struct {
	CEP string `json:"cep"`
}

type InputHandler struct {
	Client *http.Client
}

func NewInputHandler(client *http.Client) *InputHandler {
	return &InputHandler{Client: client}
}

func (h *InputHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || len(input.CEP) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	fmt.Printf("Input: %+v\n", input)

	tr := otel.Tracer("service-a")
	ctx, span := tr.Start(r.Context(), "InputHandler")
	defer span.End()

	serviceBURL := os.Getenv("SERVICE_B_URL")
	if serviceBURL == "" {
		serviceBURL = "http://service-b:8081/weather"
	}

	jsonInput, _ := json.Marshal(input)
	req, err := http.NewRequestWithContext(ctx, "POST", serviceBURL, bytes.NewBuffer(jsonInput))
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.Client.Do(req)
	if err != nil {
		http.Error(w, "failed to get response from service B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read response body", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Response from B: %s\n", string(responseBody))

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}
