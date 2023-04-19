package service

type IDoorkeeperService interface {
	HashPassword(pass string) string
	// GenerateToken(user domain.User)
}
