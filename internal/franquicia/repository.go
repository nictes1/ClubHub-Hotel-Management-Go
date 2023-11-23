package franquicia

import (
	"clubhub-hotel-management/internal/domain"
	"database/sql"
)

const (
	SAVE_USER = ""
	GET_ONE   = ""
	GET_ALL   = ""
	EXIST     = ""
	DELETE    = ""
)

type Repository interface {
	Create(f domain.Franquicia) error
	Update(f domain.Franquicia) error
	GetOne(f domain.Franquicia) error
	GetAll(f domain.Franquicia) error
}

type repo struct {
	db *sql.DB
}

// NewRepository crea un nuevo repositorio de franquicia.
func NewRepository(mysql *sql.DB) Repository {
	return &repo{
		db: mysql,
	}
}

func (r *repo) Create(f domain.Franquicia) error {
	return nil
}

func (r *repo) Update(f domain.Franquicia) error {
	return nil
}

func (r *repo) GetOne(f domain.Franquicia) error {
	return nil
}

func (r *repo) GetAll(f domain.Franquicia) error {
	return nil
}
