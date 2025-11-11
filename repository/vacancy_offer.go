package repository

import (
	"avileads-web/models"
	"avileads-web/repository/audit"
	repository "avileads-web/repository/specification"
	"github.com/astaxie/beego/orm"
)

type VacancyOfferRepository struct {
	*BaseRepository[models.Vacancy_2_Offer]
}

func NewVacancyOfferRepository(db orm.Ormer, meta *audit.AuditMeta) *VacancyOfferRepository {
	return &VacancyOfferRepository{New[models.Vacancy_2_Offer](db, meta)}
}

func (r *VacancyOfferRepository) GetAllByVacancyId(vacancyId int) ([]models.Vacancy_2_Offer, error) {
	return r.GetAll(repository.Vacancy2OffersByVacancyId{VacancyId: vacancyId})
}

func (r *VacancyOfferRepository) ByOfferId(offerId int) (*models.Vacancy_2_Offer, error) {
	return r.First(repository.Vacancy2OffersByOfferId{OfferId: offerId})
}
