package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
)

type ProfessionService struct {
	repository *repository.RkDictProfRepository
}

func NewProfessionService(r *repository.RkDictProfRepository) *ProfessionService {
	return &ProfessionService{repository: r}
}

func (s *ProfessionService) Get(id int) (*models.RkProfession, error) {
	return s.repository.Get(id)
}

func (s *ProfessionService) GetAll() ([]models.RkProfession, error) {
	return s.repository.GetAll(nil, "name")
}

func (s *ProfessionService) Create(model *models.RkProfession) (int64, error) {
	return s.repository.Create(model)
}

func (s *ProfessionService) Update(model *models.RkProfession) error {
	return s.repository.Update(model)
}

func (s *ProfessionService) Delete(id int) error {
	return s.repository.Delete(id)
}
