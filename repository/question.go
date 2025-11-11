package repository

import (
	"avileads-web/models"
	"avileads-web/repository/audit"
	repository "avileads-web/repository/specification"
	"github.com/astaxie/beego/orm"
)

type QuestionRepository struct {
	*BaseRepository[models.Questions]
}

func NewQuestionRepository(db orm.Ormer, meta *audit.AuditMeta) *QuestionRepository {
	return &QuestionRepository{New[models.Questions](db, meta)}
}

func (r *QuestionRepository) GetAllByQuestionnaireId(QuestionnaireId int) ([]models.Questions, error) {
	return r.GetAll(
		repository.QuestionsByQuestionId{QuestionnaireId: QuestionnaireId},
		"Sort",
	)
}
