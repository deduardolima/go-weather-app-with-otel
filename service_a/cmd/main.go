package main

import (
	"log"
	"net/http"
	"os"

	"github.com/deduardolima/go-weather-with-otel/internal"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans",
		zipkin.WithLogger(log.Default()),
	)
	if err != nil {
		log.Fatalf("failed to create Zipkin exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)

	tracer := otel.Tracer("service-a")

	r := mux.NewRouter()
	r.HandleFunc("/cep", internal.NewCEPHandler(tracer)).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Service A is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, otelhttp.NewHandler(r, "service-a-handler")))
}
