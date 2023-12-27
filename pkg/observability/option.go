package observability

type Option func(*observability)

func TraceEndpoint(endpoint string) Option {
	return func(o *observability) {
		o.traceExporterEndpoint = endpoint
	}
}

func MetricsEndpoint(endpoint string) Option {
	return func(o *observability) {
		o.metricsExporterEndpoint = endpoint
	}
}

func ServiceName(name string) Option {
	return func(o *observability) {
		o.serviceName = name
	}
}

func ServiceVersion(version string) Option {
	return func(o *observability) {
		o.serviceVersion = version
	}
}

func ServiceInstanceID(id string) Option {
	return func(o *observability) {
		o.serviceInstanceID = id
	}
}
