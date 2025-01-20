package repository

type Records interface {

}

type Repository struct {
	Records
}

func NewRepository() *Repository {
	return &Repository{
	}
}