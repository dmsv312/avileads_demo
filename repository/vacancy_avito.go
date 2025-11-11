package repository

import (
	"avileads-web/models"
	"avileads-web/repository/audit"
	"github.com/astaxie/beego/orm"
)

type VacancyRepository struct {
	*BaseRepository[models.Vacancy_avito]
}

func NewVacancyRepository(db orm.Ormer, meta *audit.AuditMeta) *VacancyRepository {
	return &VacancyRepository{New[models.Vacancy_avito](db, meta)}
}
