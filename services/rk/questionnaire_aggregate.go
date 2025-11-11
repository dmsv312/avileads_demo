package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
	"avileads-web/transport"
)

type QuestionDiff struct {
	ToCreate []transport.QuestionPayload
	ToUpdate []transport.QuestionPayload
	ToDelete []int
}

type QuestionnaireService struct {
	questionnaireRepository *repository.QuestionnaireRepository
	questionRepository      *repository.QuestionRepository
}

func NewQuestionnaireService(questionnaireRepository *repository.QuestionnaireRepository, questionRepository *repository.QuestionRepository) *QuestionnaireService {
	return &QuestionnaireService{
		questionnaireRepository: questionnaireRepository,
		questionRepository:      questionRepository,
	}
}

func (s *QuestionnaireService) CreateOrUpdate(dto transport.QuestionnairePayload) (int64, error) {
	if dto.Id == 0 {
		return s.create(dto)
	}
	return s.update(dto)
}

func (s *QuestionnaireService) create(dto transport.QuestionnairePayload) (int64, error) {
	questionnaire := s.dtoToModel(&models.Questionnaire{}, dto)

	id64, err := s.questionnaireRepository.Create(questionnaire)
	if err != nil {
		return 0, err
	}
	questionnaireId := int(id64)

	for _, p := range dto.Questions {
		question := s.questionDtoToModel(&models.Questions{QuestID: &models.Questionnaire{Id: questionnaireId}}, p)
		if _, err = s.questionRepository.Create(question); err != nil {
			return 0, err
		}
	}
	return id64, nil
}

func (s *QuestionnaireService) update(dto transport.QuestionnairePayload) (int64, error) {
	existing, err := s.questionRepository.GetAllByQuestionnaireId(dto.Id)
	if err != nil {
		return 0, err
	}

	diff := s.diffQuestions(existing, dto.Questions)

	if err = s.updateQuestions(diff.ToUpdate); err != nil {
		return 0, err
	}
	if err = s.createQuestions(dto.Id, diff.ToCreate); err != nil {
		return 0, err
	}
	if _, err = s.delete(diff.ToDelete); err != nil {
		return 0, err
	}

	return int64(dto.Id), nil
}

func (s *QuestionnaireService) createQuestions(questionnaireId int, payload []transport.QuestionPayload) error {
	for _, questionDTO := range payload {
		question := s.questionDtoToModel(&models.Questions{QuestID: &models.Questionnaire{Id: questionnaireId}}, questionDTO)
		if _, err := s.questionRepository.Create(question); err != nil {
			return err
		}
	}
	return nil
}

func (s *QuestionnaireService) updateQuestions(payload []transport.QuestionPayload) error {
	for _, questionDTO := range payload {
		question, err := s.questionRepository.Get(questionDTO.Id)
		if err != nil {
			return err
		}
		question = s.questionDtoToModel(question, questionDTO)
		if err = s.questionRepository.Update(question); err != nil {
			return err
		}
	}
	return nil
}

func (s *QuestionnaireService) delete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return s.questionRepository.DeleteBatch(ids)
}

func (s *QuestionnaireService) diffQuestions(
	existing []models.Questions,
	payload []transport.QuestionPayload,
) QuestionDiff {

	existSet := make(map[int]struct{}, len(existing))
	for _, q := range existing {
		existSet[q.Id] = struct{}{}
	}

	var diff QuestionDiff
	arrived := make(map[int]struct{}, len(payload))

	for _, p := range payload {
		if p.Id > 0 {
			diff.ToUpdate = append(diff.ToUpdate, p)
			arrived[p.Id] = struct{}{}
			continue
		}
		diff.ToCreate = append(diff.ToCreate, p)
	}

	for id := range existSet {
		if _, ok := arrived[id]; !ok {
			diff.ToDelete = append(diff.ToDelete, id)
		}
	}

	return diff
}

func (s *QuestionnaireService) dtoToModel(questionnaire *models.Questionnaire, questionnaireDTO transport.QuestionnairePayload) *models.Questionnaire {
	questionnaire.Name = questionnaireDTO.Name
	questionnaire.ClientId = questionnaireDTO.ClientId
	return questionnaire
}

func (s *QuestionnaireService) questionDtoToModel(question *models.Questions, questionDTO transport.QuestionPayload) *models.Questions {
	question.Text = questionDTO.Text
	question.WrongAnswerMessages = questionDTO.WrongAnswer
	question.FollowUpMessages = questionDTO.FollowUp
	question.TypeID = &models.QuestType{Id: questionDTO.TypeId}
	question.IsRequired = questionDTO.IsRequired
	question.Sort = questionDTO.Sort
	return question
}

func (s *QuestionnaireService) Get(id int) (*models.Questionnaire, error) {
	return s.questionnaireRepository.Get(id)
}

func (s *QuestionnaireService) GetDTO(id int) (transport.QuestionnaireResponse, error) {
	var dto transport.QuestionnaireResponse

	questionnaire, err := s.questionnaireRepository.Get(id)
	if err != nil {
		return dto, err
	}
	dto = s.questionnaireModelToDTO(questionnaire)

	questions, err := s.questionRepository.GetAllByQuestionnaireId(id)
	if err != nil {
		return dto, err
	}

	for _, q := range questions {
		dto.Questions = append(dto.Questions, s.questionModelToDTO(q))
	}
	return dto, nil
}

func (s *QuestionnaireService) questionnaireModelToDTO(questionnaire *models.Questionnaire) transport.QuestionnaireResponse {
	return transport.QuestionnaireResponse{
		Id:   questionnaire.Id,
		Name: questionnaire.Name,
	}
}

func (s *QuestionnaireService) questionModelToDTO(question models.Questions) transport.QuestionResponse {
	return transport.QuestionResponse{
		Id:                  question.Id,
		Text:                question.Text,
		WrongAnswerMessages: question.WrongAnswerMessages,
		FollowUpMessages:    question.FollowUpMessages,
		TypeId:              question.TypeID.Id,
		IsRequired:          question.IsRequired,
		Sort:                question.Sort,
	}
}
