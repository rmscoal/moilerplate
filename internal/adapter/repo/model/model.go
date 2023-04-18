package model

// GetAllModels returns all relational models which
// will be used for migrations needs. Register the
// models here.
func GetAllRelationalModels() []any {
	return []any{
		&User{},
		&UserCredential{},
		&UserEmail{},
	}
}
