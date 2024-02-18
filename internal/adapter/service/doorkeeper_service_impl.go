package service

import (
	"context"
	"crypto/sha1"
	"encoding/base64"

	"github.com/rmscoal/moilerplate/pkg/doorkeeper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/pbkdf2"
)

type doorkeeperService struct {
	dk     *doorkeeper.Doorkeeper
	tracer trace.Tracer
}

func NewDoorkeeperService(dk *doorkeeper.Doorkeeper) *doorkeeperService {
	return &doorkeeperService{dk: dk, tracer: otel.Tracer("doorkeeper-service")}
}

func (service *doorkeeperService) HashAndEncodeStringWithSalt(ctx context.Context, str, slt string) string {
	_, span := service.tracer.Start(ctx, "service.HashAndEncodeStringWithSalt")
	span.End()

	salt := sha1.Sum([]byte(slt))

	return base64.StdEncoding.EncodeToString(
		pbkdf2.Key(
			[]byte(str),
			salt[:],
			service.dk.GetHashIter(),
			service.dk.GetHashKeyLen(),
			service.dk.GetHasherFunc(),
		),
	)
}
