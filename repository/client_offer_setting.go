package repository

import (
	"avileads-web/models"
	"avileads-web/repository/audit"
	repository "avileads-web/repository/specification"
	"github.com/astaxie/beego/orm"
)

type ClientOfferSettingRepository struct {
	*BaseRepository[models.ClientOfferSetting]
}

func NewClientOfferSettingRepository(db orm.Ormer, meta *audit.AuditMeta) *ClientOfferSettingRepository {
	return &ClientOfferSettingRepository{New[models.ClientOfferSetting](db, meta)}
}

func (r *ClientOfferSettingRepository) ByOfferId(offerId int) (*models.ClientOfferSetting, error) {
	return r.First(repository.ClientOfferByOfferId{OfferId: offerId})
}
