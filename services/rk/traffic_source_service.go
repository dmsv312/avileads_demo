package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
)

type TrafficSourceService struct {
	repository *repository.RkDictTrafficSourceRepository
}

func NewTrafficSourceService(repository *repository.RkDictTrafficSourceRepository) *TrafficSourceService {
	return &TrafficSourceService{repository: repository}
}

func (s *TrafficSourceService) Get(id int) (*models.RkTrafficSource, error) {
	return s.repository.Get(id)
}

func (s *TrafficSourceService) GetAll() ([]models.RkTrafficSource, error) {
	return s.repository.GetAll(nil, "name")
}

func (s *TrafficSourceService) Create(model *models.RkTrafficSource) (int64, error) {
	return s.repository.Create(model)
}

func (s *TrafficSourceService) Update(model *models.RkTrafficSource) error {
	return s.repository.Update(model)
}

func (s *TrafficSourceService) Delete(id int) error {
	return s.repository.Delete(id)
}
