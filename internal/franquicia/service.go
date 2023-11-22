package franquicia

import "clubhub-hotel-management/internal/domain"

type service struct {
	repo Repository
}

type Service interface {
	CreateFranquicia(f domain.Franquicia) error
	GetFranquiciaByID(f domain.Franquicia) error
	GetAllFranquicias(f domain.Franquicia) error
	UpdateFranquicia(f domain.Franquicia) error
}

// NewService crea un nuevo servicio de franquicia.
func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateFranquicia(f domain.Franquicia) error {
	return s.repo.Create(f)
}

func (s *service) GetFranquiciaByID(f domain.Franquicia) error {
	return s.repo.GetOne(f)
}

func (s *service) GetAllFranquicias(f domain.Franquicia) error {
	return s.repo.Create(f)
}

func (s *service) UpdateFranquicia(f domain.Franquicia) error {
	return s.repo.Create(f)
}
