package repository

import (
	"avileads-web/models"
	"avileads-web/repository/audit"
	repository "avileads-web/repository/specification"
	"github.com/astaxie/beego/orm"
)

type OfferFilterRepository struct {
	*BaseRepository[models.OfferFilter]
}

func NewOfferFilterRepository(db orm.Ormer, meta *audit.AuditMeta) *OfferFilterRepository {
	return &OfferFilterRepository{New[models.OfferFilter](db, meta)}
}

func (r *OfferFilterRepository) GetAllByOfferId(offerId int) ([]models.OfferFilter, error) {
	return r.GetAll(repository.ClientOfferByOfferId{OfferId: offerId})
}
