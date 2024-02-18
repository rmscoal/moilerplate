package repo

type credentialRepo struct {
	*baseRepo
}

func NewCredentialRepo() *credentialRepo {
	return &credentialRepo{baseRepo: gormRepo}
}
