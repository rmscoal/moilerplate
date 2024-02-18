package service

import "context"

type IDoorkeeperService interface {
	HashAndEncodeStringWithSalt(ctx context.Context, str, slt string) string
}
