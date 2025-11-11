package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
	"avileads-web/transport"
)

type ClientOfferSettingService struct {
	repository *repository.ClientOfferSettingRepository
}

func NewClientOfferSettingService(r *repository.ClientOfferSettingRepository) *ClientOfferSettingService {
	return &ClientOfferSettingService{repository: r}
}

func (s *ClientOfferSettingService) Get(id int) (*models.ClientOfferSetting, error) {
	return s.repository.Get(id)
}

func (s *ClientOfferSettingService) Create(clientOfferSetting *models.ClientOfferSetting) (int64, error) {
	return s.repository.Create(clientOfferSetting)
}

func (s *ClientOfferSettingService) Update(clientOfferSetting *models.ClientOfferSetting) error {
	return s.repository.Update(clientOfferSetting)
}

func (s *ClientOfferSettingService) UpdateAll(updatedOffers []transport.OfferPayload) error {
	for _, offerPayload := range updatedOffers {
		setting, err := s.repository.ByOfferId(offerPayload.Id)
		if err != nil {
			return err
		}

		setting = s.DtoToModel(setting, offerPayload)

		if err = s.repository.Update(setting, "link", "sheet_id", "sheet_name"); err != nil {
			return err
		}
	}

	return nil
}

func (s *ClientOfferSettingService) Delete(id int) error {
	return s.repository.Delete(id)
}

func (s *ClientOfferSettingService) DeleteByOfferIds(offerIds []int) error {
	for _, oid := range offerIds {
		setting, err := s.repository.ByOfferId(oid)
		if err != nil {
			return err
		}

		err = s.repository.Delete(setting.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ClientOfferSettingService) GetDTOByOfferId(offerId int) (transport.ClientOfferSettingResponse, error) {
	clientOfferSetting, err := s.repository.ByOfferId(offerId)
	if err != nil {
		return transport.ClientOfferSettingResponse{}, err
	}

	return s.ModelToDTO(clientOfferSetting), nil
}

func (s *ClientOfferSettingService) DtoToModel(clientOfferSetting *models.ClientOfferSetting, offerDTO transport.OfferPayload) *models.ClientOfferSetting {
	clientOfferSetting.Link = offerDTO.Link
	clientOfferSetting.SheetId = offerDTO.SheetId
	clientOfferSetting.SheetName = offerDTO.SheetName

	return clientOfferSetting
}

func (s *ClientOfferSettingService) ModelToDTO(clientOfferSetting *models.ClientOfferSetting) transport.ClientOfferSettingResponse {
	return transport.ClientOfferSettingResponse{
		Id:        clientOfferSetting.Id,
		OfferId:   clientOfferSetting.OfferId,
		ClientId:  clientOfferSetting.ClientId,
		Link:      clientOfferSetting.Link,
		SheetId:   clientOfferSetting.SheetId,
		SheetName: clientOfferSetting.SheetName,
	}
}
