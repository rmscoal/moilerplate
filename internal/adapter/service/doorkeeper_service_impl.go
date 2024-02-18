package service

import (
	"github.com/rmscoal/moilerplate/pkg/doorkeeper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type doorkeeperService struct {
	dk     *doorkeeper.Doorkeeper
	tracer trace.Tracer
}

func NewDoorkeeperService(dk *doorkeeper.Doorkeeper) *doorkeeperService {
	return &doorkeeperService{dk: dk, tracer: otel.Tracer("doorkeeper-service")}
}
