package doorkeeper

import (
	"fmt"

	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/doorkeeper"
)

type doorkeeperService struct {
	dk *doorkeeper.Doorkeeper
}

func NewDoorkeeperService(dk *doorkeeper.Doorkeeper) *doorkeeperService {
	return &doorkeeperService{dk}
}

func (s *doorkeeperService) HashPassword(pass string) string {
	h := s.dk.GetHasMethod().New()
	h.Write([]byte(pass))
	res := h.Sum([]byte(s.dk.GetSalt()))

	return fmt.Sprintf("%x", res)
}
