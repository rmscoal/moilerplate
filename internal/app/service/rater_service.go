package service

import "context"

type IRaterService interface {
	IsClientAllowed(ctx context.Context, ip string) bool
}
