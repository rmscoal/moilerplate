package repo

type IBaseRepo interface {
	TranslateError(err error) error
}
