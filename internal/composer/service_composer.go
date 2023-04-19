package composer

import (
	impl "github.com/rmscoal/go-restful-monolith-boilerplate/internal/adapter/service/doorkeeper"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/service"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/doorkeeper"
)

type IServiceComposer interface {
	DoorkeeperService() service.IDoorkeeperService
}

type serviceComposer struct {
	dk *doorkeeper.Doorkeeper
}

func NewServiceComposer(dk *doorkeeper.Doorkeeper) IServiceComposer {
	return &serviceComposer{dk}
}

func (s *serviceComposer) DoorkeeperService() service.IDoorkeeperService {
	return impl.NewDoorkeeperService(s.dk)
}
