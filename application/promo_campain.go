package application

import (
	"avileads-web/repository"
	"avileads-web/repository/audit"
	"avileads-web/services"
	"avileads-web/services/rk"
	"avileads-web/transport"
	"github.com/astaxie/beego/orm"
)

type PromoCampaign struct {
	db                           orm.Ormer
	meta                         *audit.AuditMeta
	vacancyRepository            *repository.VacancyRepository
	openAIRepository             *repository.VacancyOpenAIRepository
	offerRepository              *repository.OfferRepository
	vacancyOfferRepository       *repository.VacancyOfferRepository
	offerFilterRepository        *repository.OfferFilterRepository
	clientOfferSettingRepository *repository.ClientOfferSettingRepository
	questionnaireRepository      *repository.QuestionnaireRepository
	questionRepository           *repository.QuestionRepository
	vacancyService               *rk.VacancyService
	offerService                 *rk.OfferService
	vacancyOfferService          *rk.VacancyOfferService
	offerFilterService           *rk.OfferFilterService
	clientOfferSettingService    *rk.ClientOfferSettingService
	offerAggregateService        *rk.OfferAggregateService
	questionnaireService         *rk.QuestionnaireService
}

func NewPromoCampaign(db orm.Ormer, clientId int, userId int) *PromoCampaign {
	meta := &audit.AuditMeta{UserId: userId, ClientId: clientId}

	vacancyRepository := repository.NewVacancyRepository(db, meta)
	openAIRepository := repository.NewVacancyOpenAIRepository(db, meta)
	offerRepository := repository.NewOfferRepository(db, meta)
	vacancyOfferRepository := repository.NewVacancyOfferRepository(db, meta)
	offerFilterRepository := repository.NewOfferFilterRepository(db, meta)
	clientOfferSettingRepository := repository.NewClientOfferSettingRepository(db, meta)
	questionnaireRepository := repository.NewQuestionnaireRepository(db, meta)
	questionRepository := repository.NewQuestionRepository(db, meta)

	vacancyService := rk.NewVacancyService(vacancyRepository, openAIRepository)
	offerService := rk.NewOfferService(offerRepository, vacancyOfferRepository)
	vacancyOfferService := rk.NewVacancyOfferService(vacancyOfferRepository)
	offerFilterService := rk.NewOfferFilterService(offerFilterRepository)
	clientOfferSettingService := rk.NewClientOfferSettingService(clientOfferSettingRepository)
	questionnaireService := rk.NewQuestionnaireService(questionnaireRepository, questionRepository)

	offerAggregateService := &rk.OfferAggregateService{
		OfferService:         offerService,
		FilterService:        offerFilterService,
		ClientSettingService: clientOfferSettingService,
		VacancyOfferService:  vacancyOfferService,
	}

	return &PromoCampaign{
		db:                           db,
		meta:                         meta,
		vacancyRepository:            vacancyRepository,
		offerRepository:              offerRepository,
		vacancyOfferRepository:       vacancyOfferRepository,
		offerFilterRepository:        offerFilterRepository,
		clientOfferSettingRepository: clientOfferSettingRepository,
		questionnaireRepository:      questionnaireRepository,
		questionRepository:           questionRepository,
		vacancyService:               vacancyService,
		offerService:                 offerService,
		vacancyOfferService:          vacancyOfferService,
		offerFilterService:           offerFilterService,
		clientOfferSettingService:    clientOfferSettingService,
		offerAggregateService:        offerAggregateService,
		questionnaireService:         questionnaireService,
	}
}

func (promoCampaign *PromoCampaign) Process(request transport.PromoCampaignPayload) (int, error) {
	promoCampaign.db.Begin()

	if len(request.NewDictionaries) > 0 {
		dictionaryService := rk.NewDictAggregateService(promoCampaign.db, promoCampaign.meta)
		ids, err := dictionaryService.CreateBatch(request.NewDictionaries)
		if err != nil {
			promoCampaign.db.Rollback()
			return 0, err
		}
	}

	if request.Vacancy.Id != 0 {
		promoCampaign.meta.RkId = request.Vacancy.Id
	}
	request.Vacancy.ClientId = promoCampaign.meta.ClientId

	vacancyId, err := promoCampaign.vacancyService.CreateOrUpdate(request.Vacancy)
	if err != nil {
		promoCampaign.db.Rollback()
		return 0, err
	}

	promoCampaign.meta.RkId = int(vacancyId)

	vacancy, err := promoCampaign.vacancyService.Get(int(vacancyId))
	if err != nil {
		promoCampaign.db.Rollback()
		return int(vacancyId), err
	}

	for i, _ := range request.Offers {
		request.Offers[i].ClientId = promoCampaign.meta.ClientId
		request.Offers[i].Name = vacancy.Name
	}

	err = promoCampaign.offerAggregateService.Sync(int(vacancyId), request.Offers)
	if err != nil {
		promoCampaign.db.Rollback()
		return int(vacancyId), err
	}

	if !request.Vacancy.OpenAiSupport {
		if request.Questionnaire.Id == 0 {
			request.Questionnaire.Name = vacancy.Name
		}
		request.Questionnaire.ClientId = request.Vacancy.ClientId
		questionnaireId, err := promoCampaign.questionnaireService.CreateOrUpdate(request.Questionnaire)
		if err != nil {
			promoCampaign.db.Rollback()
			return int(vacancyId), err
		}

		newVacancyDTO := request.Vacancy
		newVacancyDTO.Id = int(vacancyId)
		newVacancyDTO.Questionnaire = int(questionnaireId)

		err = promoCampaign.vacancyService.Update(newVacancyDTO)
		if err != nil {
			promoCampaign.db.Rollback()
			return int(vacancyId), err
		}
	}

	return int(vacancyId), promoCampaign.db.Commit()
}

func (promoCampaign *PromoCampaign) Get(vacancyId int) (transport.PromoCampaignResponse, error) {
	promoCampaign.db.Begin()

	response := transport.PromoCampaignResponse{}
	vacancyDTO, err := promoCampaign.vacancyService.GetDTO(vacancyId)
	if err != nil {
		promoCampaign.db.Rollback()
		return response, err
	}

	offerDTOs, err := promoCampaign.offerService.GetDTOsByVacancyId(vacancyId)
	if err != nil {
		promoCampaign.db.Rollback()
		return response, err
	}

	questionnaireDTO := transport.QuestionnaireResponse{}
	if !vacancyDTO.OpenAiSupport {
		questionnaireDTO, err = promoCampaign.questionnaireService.GetDTO(vacancyDTO.Questionnaire)
		if err != nil {
			promoCampaign.db.Rollback()
			return response, err
		}
	}

	for i := range offerDTOs {
		offerDTO := &offerDTOs[i]

		clientOfferSettingDTO, err := promoCampaign.clientOfferSettingService.GetDTOByOfferId(offerDTO.Id)
		if err != nil {
			return response, err
		}

		offerFilterDTOs, err := promoCampaign.offerFilterService.GetDTOsByOfferId(offerDTO.Id)
		if err != nil {
			return response, err
		}

		offerDTO.ClientOfferSetting = clientOfferSettingDTO
		offerDTO.OfferFilters = offerFilterDTOs
	}

	response.Vacancy = vacancyDTO
	response.Offers = offerDTOs
	response.Questionnaire = questionnaireDTO

	return response, promoCampaign.db.Commit()
}
