package rk

import (
	"avileads-web/models"
	"avileads-web/repository"
	"avileads-web/repository/audit"
	"avileads-web/transport"
	"fmt"
	"github.com/astaxie/beego/orm"
)

type DictAggregateService struct {
	BuyingTypeService    *BuyingTypeService
	ProfessionService    *ProfessionService
	ProviderService      *ProviderService
	ReklService          *ReklService
	TrafficSourceService *TrafficSourceService
}

func NewDictAggregateService(db orm.Ormer, meta *audit.AuditMeta) *DictAggregateService {
	buyingTypeRepository := repository.NewRkDictBuyingTypeRepository(db, meta)
	professionRepository := repository.NewRkDictProfRepository(db, meta)
	providerRepository := repository.NewRkDictProviderRepository(db, meta)
	reklRepository := repository.NewRkDictReklRepository(db, meta)
	trafficSourceRepository := repository.NewRkDictTrafficSourceRepository(db, meta)

	buyingTypeService := NewBuyingTypeService(buyingTypeRepository)
	professionService := NewProfessionService(professionRepository)
	providerService := NewProviderService(providerRepository)
	reklService := NewReklService(reklRepository)
	trafficSourceService := NewTrafficSourceService(trafficSourceRepository)

	return &DictAggregateService{
		BuyingTypeService:    buyingTypeService,
		ProfessionService:    professionService,
		ProviderService:      providerService,
		ReklService:          reklService,
		TrafficSourceService: trafficSourceService,
	}
}

func (a *DictAggregateService) GetBuyingTypes() ([]models.RkBuyingType, error) {
	return a.BuyingTypeService.GetAll()
}

func (a *DictAggregateService) GetProfessions() ([]models.RkProfession, error) {
	return a.ProfessionService.GetAll()
}

func (a *DictAggregateService) GetProviders() ([]models.RkProvider, error) {
	return a.ProviderService.GetAll()
}

func (a *DictAggregateService) GetRekls() ([]models.RkRekl, error) {
	return a.ReklService.GetAll()
}

func (a *DictAggregateService) GetTrafficSources() ([]models.RkTrafficSource, error) {
	return a.TrafficSourceService.GetAll()
}

func (a *DictAggregateService) createOne(dict, name string) (int64, error) {
	switch dict {
	case "reklId":
		return a.ReklService.Create(&models.RkRekl{Name: name, ShortName: name})
	case "buyingTypeId":
		return a.BuyingTypeService.Create(&models.RkBuyingType{Name: name, ShortName: name})
	case "profId":
		return a.ProfessionService.Create(&models.RkProfession{Name: name, ShortName: name})
	case "providerId":
		return a.ProviderService.Create(&models.RkProvider{Name: name, ShortName: name})
	case "trafficSourceId":
		return a.TrafficSourceService.Create(&models.RkTrafficSource{Name: name, ShortName: name})
	default:
		return 0, fmt.Errorf("unknown dict %s", dict)
	}
}

func (a *DictAggregateService) CreateBatch(dictionariesDTO []transport.NewDictionaryPayload) (map[string]int, error) {
	result := make(map[string]int, len(dictionariesDTO))
	for _, d := range dictionariesDTO {
		id, err := a.createOne(d.Dict, d.Name)
		if err != nil {
			return nil, err
		}
		result[d.Dict] = int(id)
	}
	return result, nil
}
