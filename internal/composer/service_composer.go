package composer

import (
	impl "github.com/rmscoal/moilerplate/internal/adapter/service"
	"github.com/rmscoal/moilerplate/internal/app/service"
	"github.com/rmscoal/moilerplate/pkg/doorkeeper"
	"github.com/rmscoal/moilerplate/pkg/rater"
)

type IServiceComposer interface {
	DoorkeeperService() service.IDoorkeeperService
	RaterService() service.IRaterService
}

type serviceComposer struct {
	dk *doorkeeper.Doorkeeper
	rt *rater.Rater
}

func NewServiceComposer(dk *doorkeeper.Doorkeeper, rt *rater.Rater) IServiceComposer {
	return &serviceComposer{dk: dk, rt: rt}
}

func (s *serviceComposer) DoorkeeperService() service.IDoorkeeperService {
	return impl.NewDoorkeeperService(s.dk)
}

func (s *serviceComposer) RaterService() service.IRaterService {
	return impl.NewRaterService(s.rt)
}
