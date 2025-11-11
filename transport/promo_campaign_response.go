package transport

type PromoCampaignResponse struct {
	Vacancy       VacancyResponse       `json:"vacancy"`
	Offers        []OfferResponse       `json:"offers"`
	Questionnaire QuestionnaireResponse `json:"questionnaire"`
}

type VacancyResponse struct {
	Id                                 int            `json:"id"`
	Name                               string         `json:"name"`
	LoadToZp                           bool           `json:"loadToZp"`
	Questionnaire                      int            `json:"questionnaire"`
	HelloText                          string         `json:"helloText"`
	Description                        string         `json:"Description"`
	ByeText                            string         `json:"byeText"`
	EnablePostProcessing               bool           `json:"enablePostProcessing"`
	PostProcessingMessages             string         `json:"postProcessingMessages"`
	DialogLifeTimeInMinutes            int            `json:"dialogLifeTimeInMinutes"`
	FollowUpMessageIntervalInMinutes   int            `json:"followUpMessageIntervalInMinutes"`
	CountOfMessagesAfterFinishedDialog *int64         `json:"countOfMessagesAfterFinishedDialog,omitempty"`
	OpenAiSupport                      bool           `json:"column(openAiSupport)"`
	OpenAIDetails                      OpenAIResponse `json:"openAIDetails,omitempty"`

	ReklId             int    `json:"reklId"`
	BuyingTypeId       int    `json:"buyingTypeId"`
	ProfessionId       int    `json:"profId"`
	ProviderId         int    `json:"providerId"`
	TrafficSourceId    int    `json:"trafficSourceId"`
	ExportTypeId       int    `json:"exportTypeId"`
	CpaExId            string `json:"cpaExId,omitempty"`
	OtherSpecification string `json:"otherSpecification,omitempty"`
}

type OpenAIResponse struct {
	Id                   int     `json:"id"`
	VacancyDescription   string  `json:"vacancyDescription"`
	AssistantDescription string  `json:"assistantDescription"`
	Questions            string  `json:"questions"`
	AssistantTemperature float32 `json:"assistantTemperature"`
}

type OfferResponse struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	SqlForZp     string `json:"sqlForZp"`
	ZpScriptName string `json:"zpScriptName"`
	Enable       bool   `json:"enable"`
	IgnoreName   bool   `json:"ignoreName"`
	ExportZpPg   bool   `json:"exportZpPg"`
	ExportTypeId int    `json:"exportTypeId"`
	ClientId     int    `json:"clientId"`

	ClientOfferSetting ClientOfferSettingResponse `json:"clientOfferSetting,omitempty"`
	OfferFilters       []OfferFilterResponse      `json:"offerFilters"`
}

type SQLFilterResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type OfferFilterResponse struct {
	Id          int `json:"id"`
	OfferId     int `json:"offerId"`
	SqlFilterId int `json:"sqlFilterId"`
	ClientId    int `json:"clientId"`
}

type ClientOfferSettingResponse struct {
	Id        int    `json:"id"`
	OfferId   int    `json:"offerId"`
	ClientId  int    `json:"clientId"`
	Link      string `json:"link"`
	SheetId   string `json:"sheetId"`
	SheetName string `json:"sheetName"`
}

type ExportTypeResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ClientResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type QuestionnaireResponse struct {
	Id        int                `json:"id"`
	Name      string             `json:"name"`
	Questions []QuestionResponse `json:"questions"`
}

type QuestionResponse struct {
	Id                  int    `json:"id"`
	Text                string `json:"text"`
	WrongAnswerMessages string `json:"wrongAnswerMessage"`
	FollowUpMessages    string `json:"followUp"`
	TypeId              int    `json:"typeId"`
	IsRequired          bool   `json:"isRequired"`
	Sort                int    `json:"sort"`
}

type RkDictBuyingTypeResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

type RkDictProfResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	BasicName string `json:"basicName"`
}

type RkDictProviderResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Type      string `json:"type"`
}

type RkDictReklResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

type RkDictTrafficSourceResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}
