package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
	"avileads-web/transport"
	"database/sql"
	"strings"

	"github.com/astaxie/beego/orm"
)

type VacancyService struct {
	repository       *repository.VacancyRepository
	openAIRepository *repository.VacancyOpenAIRepository
}

func NewVacancyService(vacancyRepository *repository.VacancyRepository, openAIRepository *repository.VacancyOpenAIRepository) *VacancyService {
	return &VacancyService{repository: vacancyRepository, openAIRepository: openAIRepository}
}

func (s *VacancyService) Get(id int) (*models.Vacancy_avito, error) {
	return s.repository.Get(id)
}

func (s *VacancyService) GetDTO(id int) (transport.VacancyResponse, error) {
	vacancy, err := s.repository.Get(id)
	if err != nil {
		return transport.VacancyResponse{}, err
	}

	vacancyDTO := s.ModelToDTO(vacancy)

	if vacancy.OpenAI_Support {
		openAIDetails, err := s.openAIRepository.ByVacancyId(vacancy.Id)
		if err != nil {
			return s.ModelToDTO(vacancy), nil
		}
		vacancyDTO.OpenAIDetails = s.OpenAIModelToDTO(openAIDetails)
	}

	return vacancyDTO, nil
}

func (s *VacancyService) Create(vacancy *models.Vacancy_avito) (int64, error) {
	id, err := s.repository.Create(vacancy)
	if err != nil {
		return 0, err
	}

	fullVacancy, err := s.repository.Get(int(id),
		"Rekl", "Profession", "Provider",
		"ExportType", "BuyingType", "TrafficSource",
	)
	if err != nil {
		return 0, err
	}

	fullVacancy.Name = buildVacancyName(fullVacancy)
	if err = s.repository.Update(fullVacancy, "Name"); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *VacancyService) Update(vacancyDTO transport.VacancyPayload) error {
	vacancy, err := s.Get(vacancyDTO.Id)
	if err != nil {
		return err
	}

	vacancy = s.DtoToModel(vacancy, vacancyDTO)
	err = s.repository.Update(vacancy)
	if err != nil {
		return err
	}

	fullVacancy, err := s.repository.Get(int(vacancy.Id),
		"Rekl", "Profession", "Provider",
		"ExportType", "BuyingType", "TrafficSource",
	)
	if err != nil {
		return err
	}

	fullVacancy.Name = buildVacancyName(fullVacancy)
	if err = s.repository.Update(fullVacancy, "Name"); err != nil {
		return err
	}

	return nil
}

func (s *VacancyService) CreateOrUpdate(vacancyDTO transport.VacancyPayload) (int64, error) {
	if vacancyDTO.Id == 0 {
		vacancy := s.DtoToModel(&models.Vacancy_avito{}, vacancyDTO)
		id64, err := s.Create(vacancy)
		if err != nil {
			return 0, err
		}
		id := int(id64)

		if vacancyDTO.OpenAIDetails != (transport.OpenAIPayload{}) {
			vacancyDTO.OpenAIDetails.Id = 0
			openAIModel := s.OpenAIDTOToModel(
				&models.Vacancy_avito_openai{Vacancy_Avito_Id: id},
				vacancyDTO.OpenAIDetails,
			)

			if _, err = s.openAIRepository.Create(openAIModel); err != nil {
				return 0, err
			}
		}
		return id64, nil
	}

	vacancy, err := s.Get(vacancyDTO.Id)
	if err != nil {
		return 0, err
	}
	vacancy = s.DtoToModel(vacancy, vacancyDTO)
	if err = s.repository.Update(vacancy); err != nil {
		return 0, err
	}

	openAIRecord, err := s.openAIRepository.ByVacancyId(vacancy.Id)
	switch {
	case err == orm.ErrNoRows:
		if vacancyDTO.OpenAIDetails != (transport.OpenAIPayload{}) {
			openAIRecord = s.OpenAIDTOToModel(
				&models.Vacancy_avito_openai{Vacancy_Avito_Id: vacancy.Id},
				vacancyDTO.OpenAIDetails,
			)
			if _, err = s.openAIRepository.Create(openAIRecord); err != nil {
				return 0, err
			}
		}
	case err != nil:
		return 0, err
	default:
		if vacancyDTO.OpenAIDetails != (transport.OpenAIPayload{}) {
			openAIRecord = s.OpenAIDTOToModel(openAIRecord, vacancyDTO.OpenAIDetails)
			if err = s.openAIRepository.Update(openAIRecord); err != nil {
				return 0, err
			}
		}
	}

	return int64(vacancy.Id), nil
}

func (s *VacancyService) Delete(id int) error {
	return s.repository.Delete(id)
}

func (s *VacancyService) DtoToModel(vacancy *models.Vacancy_avito, vacancyDTO transport.VacancyPayload) *models.Vacancy_avito {
	vacancy.LoadToZp = vacancyDTO.LoadToZp
	vacancy.OpenAI_Support = vacancyDTO.OpenAiSupport
	vacancy.HelloText = vacancyDTO.HelloText
	vacancy.ByeText = vacancyDTO.ByeText
	vacancy.ClientId = vacancyDTO.ClientId
	vacancy.PostProcessingMessages = vacancyDTO.PostProcessingMessages
	vacancy.DialogLifeTimeInMinutes = vacancyDTO.DialogLifeTimeInMinutes

	return vacancy
}

func (s *VacancyService) ModelToDTO(vacancy *models.Vacancy_avito) transport.VacancyResponse {
	questId := 0
	if !vacancy.OpenAI_Support {
		questId = vacancy.Quest.Id
	}

	return transport.VacancyResponse{
		Id:                               vacancy.Id,
		LoadToZp:                         vacancy.LoadToZp,
		Questionnaire:                    questId,
		HelloText:                        vacancy.HelloText,
		ByeText:                          vacancy.ByeText,
		PostProcessingMessages:           vacancy.PostProcessingMessages,
		OpenAiSupport:                    vacancy.OpenAI_Support,
		DialogLifeTimeInMinutes:          vacancy.DialogLifeTimeInMinutes,
		FollowUpMessageIntervalInMinutes: vacancy.FollowUpMessageIntervalInMinutes,
	}
}

func (s *VacancyService) OpenAIDTOToModel(openAIModel *models.Vacancy_avito_openai, openAIDTO transport.OpenAIPayload) *models.Vacancy_avito_openai {
	openAIModel.Id = openAIDTO.Id
	openAIModel.VacancyDescription = openAIDTO.VacancyDescription
	openAIModel.AssistantDescription = openAIDTO.AssistantDescription
	openAIModel.Questions = openAIDTO.Questions
	openAIModel.AssistantTemperature = openAIDTO.AssistantTemperature
	return openAIModel
}

func (s *VacancyService) OpenAIModelToDTO(vacancyOpenAI *models.Vacancy_avito_openai) transport.OpenAIResponse {
	return transport.OpenAIResponse{
		Id:                   vacancyOpenAI.Id,
		VacancyDescription:   vacancyOpenAI.VacancyDescription,
		AssistantDescription: vacancyOpenAI.AssistantDescription,
		Questions:            vacancyOpenAI.Questions,
		AssistantTemperature: vacancyOpenAI.AssistantTemperature,
	}
}

func (s *VacancyService) UpdateAll() error {
	return nil
}
