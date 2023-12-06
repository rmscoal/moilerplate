package service

import (
	"context"

	"github.com/rmscoal/moilerplate/pkg/rater"
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
		client = service.rt.AddNewClient(ip)
		return service.rt.IsClientAllowed(client)
	}

	return service.rt.IsClientAllowed(client)
}
