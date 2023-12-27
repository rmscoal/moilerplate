package config

import (
	"log"
	"os"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type otelConfig struct {
	traceEndpoint  string
	metricEndpoint string

	serviceName       string
	serviceVersion    string
	serviceInstanceID string
}

func (c *Config) newOtelConfig() {
	otel := otelConfig{
		traceEndpoint:     os.Getenv("OTEL_TRACE_ENDPOINT"),
		metricEndpoint:    os.Getenv("OTEL_METRIC_ENDPOINT"),
		serviceName:       os.Getenv("OTEL_SERVICE_NAME"),
		serviceVersion:    os.Getenv("OTEL_SERVICE_VERSION"),
		serviceInstanceID: os.Getenv("OTEL_SERVICE_INSTANCE_ID"),
	}

	if err := otel.validate(); err != nil {
		log.Fatalf("FATAL - getting otel config %s\n", err)
	}

	c.Otel = otel
}

func (o otelConfig) validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.traceEndpoint, validation.Required, is.RequestURL),
		validation.Field(&o.metricEndpoint, validation.Required, is.RequestURL),
		validation.Field(&o.serviceName, validation.Required, validation.Length(3, 0)),
		validation.Field(&o.serviceVersion, validation.Required, is.Semver),
		validation.Field(&o.serviceInstanceID, validation.Required, validation.Length(3, 0)),
	)
}

func (o otelConfig) GetTraceEndpoint() string {
	return o.traceEndpoint
}

func (o otelConfig) GetMetricEndpoint() string {
	return o.metricEndpoint
}

func (o otelConfig) GetServiceName() string {
	return o.serviceName
}

func (o otelConfig) GetServiceVersion() string {
	return o.serviceVersion
}

func (o otelConfig) GetServiceInstanceID() string {
	return o.serviceInstanceID
}
