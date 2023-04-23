package service

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/rater"
)

type raterService struct {
	rt *rater.Rater
}

func NewRaterService(rt *rater.Rater) *raterService {
	return &raterService{rt}
}

func (service *raterService) IsClientAllowed(ctx context.Context, ip string) bool {
	client, found := service.rt.GetClient(ip)
	if !found {
		service.rt.AddNewClient(ip)
	}

	return service.rt.IsClientAllowed(client)
}
