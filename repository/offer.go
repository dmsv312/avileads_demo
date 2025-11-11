package repository

import (
	"avileads-web/models"
	"avileads-web/repository/audit"
	"github.com/astaxie/beego/orm"
)

type OfferRepository struct{ *BaseRepository[models.Offer] }

func NewOfferRepository(db orm.Ormer, meta *audit.AuditMeta) *OfferRepository {
	return &OfferRepository{New[models.Offer](db, meta)}
}
