package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
)

type BuyingTypeService struct {
	repository *repository.RkDictBuyingTypeRepository
}

func NewBuyingTypeService(repository *repository.RkDictBuyingTypeRepository) *BuyingTypeService {
	return &BuyingTypeService{repository: repository}
}

func (s *BuyingTypeService) Get(id int) (*models.RkBuyingType, error) {
	return s.repository.Get(id)
}

func (s *BuyingTypeService) GetAll() ([]models.RkBuyingType, error) {
	return s.repository.GetAll(nil, "name")
}

func (s *BuyingTypeService) Create(model *models.RkBuyingType) (int64, error) {
	return s.repository.Create(model)
}

func (s *BuyingTypeService) Update(model *models.RkBuyingType) error {
	return s.repository.Update(model)
}

func (s *BuyingTypeService) Delete(id int) error {
	return s.repository.Delete(id)
}
