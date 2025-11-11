package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
)

type ReklService struct {
	repository *repository.RkDictReklRepository
}

func NewReklService(repository *repository.RkDictReklRepository) *ReklService {
	return &ReklService{repository: repository}
}

func (s *ReklService) Get(id int) (*models.RkRekl, error) {
	return s.repository.Get(id)
}

func (s *ReklService) GetAll() ([]models.RkRekl, error) {
	return s.repository.GetAll(nil, "name")
}

func (s *ReklService) Create(model *models.RkRekl) (int64, error) {
	return s.repository.Create(model)
}

func (s *ReklService) Update(model *models.RkRekl) error {
	return s.repository.Update(model)
}

func (s *ReklService) Delete(id int) error {
	return s.repository.Delete(id)
}
