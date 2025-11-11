package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
	"avileads-web/transport"
)

type OfferService struct {
	offerRepository        *repository.OfferRepository
	vacancyOfferRepository *repository.VacancyOfferRepository
}

func NewOfferService(offerRepository *repository.OfferRepository, vacancyOfferRepository *repository.VacancyOfferRepository) *OfferService {
	return &OfferService{offerRepository: offerRepository, vacancyOfferRepository: vacancyOfferRepository}
}

func (s *OfferService) CreateOrUpdateAll(offers []transport.OfferPayload) error {
	for i := range offers {
		newId, err := s.CreateOrUpdate(offers[i])
		if err != nil {
			return err
		}
		offers[i].Id = newId
	}
	return nil
}

func (s *OfferService) Get(id int) (*models.Offer, error) {
	return s.offerRepository.Get(id)
}

func (s *OfferService) Create(offer *models.Offer) (int64, error) {
	return s.offerRepository.Create(offer)
}

func (s *OfferService) Update(offer *models.Offer) error {
	return s.offerRepository.Update(offer)
}

func (s *OfferService) Delete(id int) error {
	return s.offerRepository.Delete(id)
}

func (s *OfferService) DeleteBatch(ids []int) error {
	for _, id := range ids {
		err := s.Delete(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *OfferService) GetDTO(id int) (transport.OfferResponse, error) {
	offer, err := s.offerRepository.Get(id)
	if err != nil {
		return transport.OfferResponse{}, err
	}
	return s.ModelToDTO(offer), nil
}

func (s *OfferService) GetDTOsByVacancyId(vacancyId int) ([]transport.OfferResponse, error) {
	vacancy2offers, err := s.vacancyOfferRepository.GetAllByVacancyId(vacancyId)
	if err != nil {
		return []transport.OfferResponse{}, err
	}

	offerDTOs := make([]transport.OfferResponse, 0, len(vacancy2offers))

	for _, vacancy2offer := range vacancy2offers {
		offer, err := s.offerRepository.Get(vacancy2offer.Offer.Id)
		if err != nil {
			return []transport.OfferResponse{}, nil
		}
		offerDTOs = append(offerDTOs, s.ModelToDTO(offer))
	}

	return offerDTOs, nil
}

func (s *OfferService) CreateOrUpdate(offerDTO transport.OfferPayload) (int, error) {
	if offerDTO.Id == 0 {
		offer := &models.Offer{}
		offer = s.DtoToModel(offer, offerDTO)
		newId, err := s.Create(offer)
		if err != nil {
			return 0, err
		}
		return int(newId), nil
	}

	offer, err := s.offerRepository.Get(offerDTO.Id)
	if err != nil {
		return 0, err
	}

	offer = s.DtoToModel(offer, offerDTO)
	if err := s.Update(offer); err != nil {
		return 0, err
	}

	return offerDTO.Id, nil
}

func (s *OfferService) DtoToModel(offer *models.Offer, offerDTO transport.OfferPayload) *models.Offer {
	offer.Name = offerDTO.Name
	offer.SqlForZp = offerDTO.SqlForZp
	offer.ZpScriptName = offerDTO.ZpScriptName
	offer.Enable = offerDTO.Enable
	offer.IgnoreName = offerDTO.IgnoreName
	offer.ExportZpPg = offerDTO.ExportZpPg
	offer.ExportTypeId = offerDTO.ExportTypeId

	return offer
}

func (s *OfferService) ModelToDTO(offer *models.Offer) transport.OfferResponse {
	return transport.OfferResponse{
		Id:           offer.Id,
		Name:         offer.Name,
		SqlForZp:     offer.SqlForZp,
		ZpScriptName: offer.ZpScriptName,
		Enable:       offer.Enable,
		IgnoreName:   offer.IgnoreName,
		ExportZpPg:   offer.ExportZpPg,
		ExportTypeId: offer.ExportTypeId,

		ClientOfferSetting: transport.ClientOfferSettingResponse{},
		OfferFilters:       []transport.OfferFilterResponse{},
	}
}
