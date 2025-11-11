package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
)

type VacancyOfferService struct {
	repository *repository.VacancyOfferRepository
}

func NewVacancyOfferService(r *repository.VacancyOfferRepository) *VacancyOfferService {
	return &VacancyOfferService{repository: r}
}

func (s *VacancyOfferService) Get(id int) (*models.Vacancy_2_Offer, error) {
	return s.repository.Get(id)
}

func (s *VacancyOfferService) Create(vacancyOffer *models.Vacancy_2_Offer) (int64, error) {
	return s.repository.Create(vacancyOffer)
}

func (s *VacancyOfferService) Update(vacancyOffer *models.Vacancy_2_Offer) error {
	return s.repository.Update(vacancyOffer)
}

func (s *VacancyOfferService) Delete(id int) error {
	return s.repository.Delete(id)
}

func (s *VacancyOfferService) DeleteByOfferIds(offerIds []int) error {
	for _, offerId := range offerIds {
		vacancyOffer, err := s.repository.ByOfferId(offerId)
		if err != nil {
			return err
		}

		err = s.repository.Delete(vacancyOffer.Id)
		if err != nil {
			return err
		}
	}

	return nil
}
