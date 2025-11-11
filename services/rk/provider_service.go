package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
)

type ProviderService struct {
	repository *repository.RkDictProviderRepository
}

func NewProviderService(repository *repository.RkDictProviderRepository) *ProviderService {
	return &ProviderService{repository: repository}
}

func (s *ProviderService) Get(id int) (*models.RkProvider, error) {
	return s.repository.Get(id)
}

func (s *ProviderService) GetAll() ([]models.RkProvider, error) {
	return s.repository.GetAll(nil, "name")
}

func (s *ProviderService) Create(model *models.RkProvider) (int64, error) {
	return s.repository.Create(model)
}

func (s *ProviderService) Update(model *models.RkProvider) error {
	return s.repository.Update(model)
}

func (s *ProviderService) Delete(id int) error {
	return s.repository.Delete(id)
}
