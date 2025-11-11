package repository

import "github.com/astaxie/beego/orm"

type Specification interface {
	Cond() *orm.Condition
}

type ClientOfferByOfferId struct {
	OfferId int
}

func (s ClientOfferByOfferId) Cond() *orm.Condition {
	return orm.NewCondition().
		And("offer_id", s.OfferId)
}

type Vacancy2OffersByVacancyId struct {
	VacancyId int
}

func (s Vacancy2OffersByVacancyId) Cond() *orm.Condition {
	return orm.NewCondition().
		And("vacancy_avito_id", s.VacancyId)
}

type Vacancy2OffersByOfferId struct {
	OfferId int
}

func (s Vacancy2OffersByOfferId) Cond() *orm.Condition {
	return orm.NewCondition().
		And("vacancy_id", s.OfferId)
}

type VacancyOpenAIByVacancyId struct {
	VacancyId int
}

func (s VacancyOpenAIByVacancyId) Cond() *orm.Condition {
	return orm.NewCondition().
		And("vacancy_avito_id", s.VacancyId)
}

type QuestionsByQuestionId struct {
	QuestionnaireId int
}

func (s QuestionsByQuestionId) Cond() *orm.Condition {
	return orm.NewCondition().
		And("quest_id", s.QuestionnaireId)
}
