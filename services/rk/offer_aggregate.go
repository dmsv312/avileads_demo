package rk

import (
	"avileads-web/models"
	"avileads-web/transport"
	"github.com/astaxie/beego/orm"
)

type OfferDiff struct {
	ToCreate []transport.OfferPayload
	ToUpdate []transport.OfferPayload
	ToDelete []int
}

type OfferAggregateService struct {
	OfferService         *OfferService
	FilterService        *OfferFilterService
	ClientSettingService *ClientOfferSettingService
	VacancyOfferService  *VacancyOfferService
}

func (a *OfferAggregateService) Sync(vacancyID int, payload []transport.OfferPayload) error {
	existing, err := a.OfferService.GetDTOsByVacancyId(vacancyID)
	if err != nil {
		return err
	}

	diff := a.DiffOffers(existing, payload)

	if err = a.Delete(diff.ToDelete); err != nil {
		return err
	}

	if err = a.Update(diff.ToUpdate); err != nil {
		return err
	}

	if err = a.Create(vacancyID, diff.ToCreate); err != nil {
		return err
	}

	return nil
}

func (a *OfferAggregateService) Create(vacancyID int, toCreate []transport.OfferPayload) error {
	for _, dto := range toCreate {
		offer := a.OfferService.DtoToModel(&models.Offer{}, dto)
		id64, err := a.OfferService.offerRepository.Create(offer)
		if err != nil {
			return err
		}

		if _, err = a.VacancyOfferService.Create(&models.Vacancy_2_Offer{
			VacancyAvito: &models.Vacancy_avito{Id: vacancyID},
			Offer:        &models.Offer{Id: int(id64)},
		}); err != nil {
			return err
		}

		if _, err = a.ClientSettingService.Create(&models.ClientOfferSetting{
			OfferId:   int(id64),
			ClientId:  dto.ClientId,
			Link:      dto.Link,
			SheetId:   dto.SheetId,
			SheetName: dto.SheetName,
		}); err != nil {
			return err
		}

		if err = a.FilterService.CreateAll(int(id64), dto.ClientId, dto.Filters); err != nil {
			return err
		}
	}
	return nil
}

func (a *OfferAggregateService) Update(toUpdate []transport.OfferPayload) error {
	for _, dto := range toUpdate {
		offer, err := a.OfferService.offerRepository.Get(dto.Id)
		if err != nil {
			return err
		}
		offer = a.OfferService.DtoToModel(offer, dto)
		if err = a.OfferService.offerRepository.Update(offer); err != nil {
			return err
		}

		setting, err := a.ClientSettingService.repository.ByOfferId(dto.Id)
		switch {
		case err == orm.ErrNoRows:
			_, err = a.ClientSettingService.Create(a.ClientSettingService.DtoToModel(setting, dto))
		case err == nil:
			setting = a.ClientSettingService.DtoToModel(setting, dto)
			err = a.ClientSettingService.repository.Update(setting, "link", "sheet_id", "sheet_name")
		}
		if err != nil {
			return err
		}

		err = a.HandleFilters(dto)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *OfferAggregateService) Delete(toDelete []int) error {
	if len(toDelete) == 0 {
		return nil
	}
	if err := a.FilterService.DeleteByOfferIds(toDelete); err != nil {
		return err
	}
	if err := a.ClientSettingService.DeleteByOfferIds(toDelete); err != nil {
		return err
	}
	if err := a.VacancyOfferService.DeleteByOfferIds(toDelete); err != nil {
		return err
	}

	return a.OfferService.DeleteBatch(toDelete)
}

func (a *OfferAggregateService) DiffOffers(
	existing []transport.OfferResponse,
	payload []transport.OfferPayload,
) OfferDiff {

	existSet := make(map[int]struct{}, len(existing))
	for _, offerDTO := range existing {
		existSet[offerDTO.Id] = struct{}{}
	}

	var diff OfferDiff
	arrivedSet := make(map[int]struct{}, len(payload))

	for _, offerDTO := range payload {
		switch {
		case offerDTO.Id <= 0:
			diff.ToCreate = append(diff.ToCreate, offerDTO)
		default:
			arrivedSet[offerDTO.Id] = struct{}{}
			if _, ok := existSet[offerDTO.Id]; ok {
				diff.ToUpdate = append(diff.ToUpdate, offerDTO)
			}
		}
	}

	for id := range existSet {
		if _, ok := arrivedSet[id]; !ok {
			diff.ToDelete = append(diff.ToDelete, id)
		}
	}

	return diff
}

func (a *OfferAggregateService) HandleFilters(dto transport.OfferPayload) error {
	existing, err := a.FilterService.GetDTOsByOfferId(dto.Id)
	if err != nil {
		return err
	}

	existMap := make(map[int]int, len(existing))
	for _, filter := range existing {
		existMap[filter.SqlFilterId] = filter.Id
	}

	incomingIds := make(map[int]struct{}, len(dto.Filters))
	for _, filterId := range dto.Filters {
		incomingIds[filterId] = struct{}{}
	}

	var toDeleteIds []int
	for sqlFilterId, filterId := range existMap {
		if _, ok := incomingIds[sqlFilterId]; !ok {
			toDeleteIds = append(toDeleteIds, filterId)
		} else {
			delete(incomingIds, sqlFilterId)
		}
	}

	var toCreate []int
	for sqlFilterId := range incomingIds {
		toCreate = append(toCreate, sqlFilterId)
	}

	if len(toDeleteIds) > 0 {
		if err = a.FilterService.DeleteAll(toDeleteIds); err != nil {
			return err
		}
	}

	if len(toCreate) > 0 {
		if err = a.FilterService.CreateAll(dto.Id, dto.ClientId, toCreate); err != nil {
			return err
		}
	}

	return nil
}
