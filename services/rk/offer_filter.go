package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
	"avileads-web/transport"
)

type OfferFilterService struct {
	repository *repository.OfferFilterRepository
}

func NewOfferFilterService(r *repository.OfferFilterRepository) *OfferFilterService {
	return &OfferFilterService{repository: r}
}

func (s *OfferFilterService) Get(id int) (*models.OfferFilter, error) {
	return s.repository.Get(id)
}

func (s *OfferFilterService) Create(offerFilter *models.OfferFilter) (int64, error) {
	return s.repository.Create(offerFilter)
}

func (s *OfferFilterService) CreateAll(offerId, clientId int, addedFilters []int) error {
	for _, id := range addedFilters {
		_, err := s.Create(&models.OfferFilter{
			OfferId:     offerId,
			ClientId:    clientId,
			SqlFilterId: id,
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func (s *OfferFilterService) Delete(id int) error {
	return s.repository.Delete(id)
}

func (s *OfferFilterService) DeleteAll(deletedFilters []int) error {
	for _, id := range deletedFilters {
		err := s.Delete(id)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *OfferFilterService) DeleteByOfferIds(offerIds []int) error {
	for _, offerId := range offerIds {
		rows, err := s.repository.GetAllByOfferId(offerId)
		if err != nil {
			return err
		}
		var ids []int
		for _, r := range rows {
			ids = append(ids, r.Id)
		}
		if err = s.DeleteAll(ids); err != nil {
			return err
		}
	}
	return nil
}

func (s *OfferFilterService) GetDTOsByOfferId(offerId int) ([]transport.OfferFilterResponse, error) {
	offerFilters, err := s.repository.GetAllByOfferId(offerId)
	if err != nil {
		return []transport.OfferFilterResponse{}, err
	}

	filterDTOs := make([]transport.OfferFilterResponse, 0, len(offerFilters))
	for _, filter := range offerFilters {
		filterDTOs = append(filterDTOs, s.ModelToDTO(filter))
	}

	return filterDTOs, nil
}

func (s *OfferFilterService) ModelToDTO(offerFilter models.OfferFilter) transport.OfferFilterResponse {
	return transport.OfferFilterResponse{
		Id:          offerFilter.Id,
		OfferId:     offerFilter.OfferId,
		SqlFilterId: offerFilter.SqlFilterId,
		ClientId:    offerFilter.ClientId,
	}
}
