package repository

import (
	"avileads-web/models"
	"avileads-web/repository/audit"
	repository "avileads-web/repository/specification"
	"github.com/astaxie/beego/orm"
)

type VacancyOpenAIRepository struct {
	*BaseRepository[models.Vacancy_avito_openai]
}

func NewVacancyOpenAIRepository(db orm.Ormer, meta *audit.AuditMeta) *VacancyOpenAIRepository {
	return &VacancyOpenAIRepository{New[models.Vacancy_avito_openai](db, meta)}
}

func (r *VacancyOpenAIRepository) ByVacancyId(vacancyId int) (*models.Vacancy_avito_openai, error) {
	return r.First(repository.VacancyOpenAIByVacancyId{VacancyId: vacancyId})
}
