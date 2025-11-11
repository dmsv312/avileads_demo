package transport

type PromoCampaignPayload struct {
	Vacancy         VacancyPayload         `json:"vacancy"`
	Offers          []OfferPayload         `json:"offers"`
	Questionnaire   QuestionnairePayload   `json:"questionnaire"`
	NewDictionaries []NewDictionaryPayload `json:"newDictionaries,omitempty"`
	ManualAvitoIds  string                 `json:"manualAvitoIds,omitempty"`
}

type VacancyPayload struct {
	Id                                 int           `json:"id"`
	LoadToZp                           bool          `json:"loadToZP"`
	Questionnaire                      int           `json:"questionnaire"`
	HelloText                          string        `json:"helloText"`
	ByeText                            string        `json:"byeText"`
	ClientId                           int           `json:"clientId"`
	EnablePostProcessing               bool          `json:"enablePostProcessing"`
	OpenAiSupport                      bool          `json:"openAiSupport"`
	PostProcessingMessages             string        `json:"postProcessingMessages"`
	DialogLifeTimeInMinutes            int           `json:"dialogLifeTimeInMinutes"`
	FollowUpMessageIntervalInMinutes   int           `json:"followUpMessageIntervalInMinutes"`
	CountOfMessagesAfterFinishedDialog int           `json:"countOfMessagesAfterFinishedDialog"`
	OpenAIDetails                      OpenAIPayload `json:"openAIDetails,omitempty"`
	IdAvitoList                        []string      `json:"idAvitoList,omitempty"`
	ExportTypeId                       int           `json:"exportTypeId,omitempty"`
	CpaExId                            string        `json:"cpaExId,omitempty"`
	OtherSpecification                 string        `json:"otherSpecification,omitempty"`
	Description                        string        `json:"description"`

	ReklId          int `json:"reklId"`
	BuyingTypeId    int `json:"buyingTypeId"`
	ProfessionId    int `json:"profId"`
	ProviderId      int `json:"providerId"`
	TrafficSourceId int `json:"trafficSourceId"`
}

type OpenAIPayload struct {
	Id                   int     `json:"id"`
	VacancyDescription   string  `json:"vacancyDescription"`
	AssistantDescription string  `json:"assistantDescription"`
	Questions            string  `json:"questions"`
	AssistantTemperature float32 `json:"assistantTemperature"`
}

type OfferPayload struct {
	TempId       int    `json:"tempId"`
	Id           int    `json:"id"`
	Name         string `json:"name"`
	SqlForZp     string `json:"sqlForZp"`
	ZpScriptName string `json:"zpScriptName"`
	Enable       bool   `json:"enable"`
	IgnoreName   bool   `json:"ignoreName"`
	ExportZpPg   bool   `json:"exportZpPg"`
	ExportTypeId int    `json:"exportTypeId"`
	ClientId     int    `json:"clientId"`
	Link         string `json:"link"`
	SheetId      string `json:"sheetId"`
	SheetName    string `json:"sheetName"`

	Filters []int `json:"filter_ids,omitempty"`
}

type QuestionnairePayload struct {
	Id        int               `json:"id"`
	Name      string            `json:"name"`
	ClientId  int               `json:"clientId"`
	Questions []QuestionPayload `json:"questions"`
}

type QuestionPayload struct {
	Id          int    `json:"id"`
	Text        string `json:"text"`
	WrongAnswer string `json:"wrongAnswer"`
	FollowUp    string `json:"followUp"`
	TypeId      int    `json:"typeId"`
	IsRequired  bool   `json:"isRequired"`
	Sort        int    `json:"sort"`
}

type NewDictionaryPayload struct {
	Dict string `json:"dict"`
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CheckDocumentPayload struct {
	InputType      string `form:"inputType"`
	ManualAvitoIds string `form:"manualAvitoIds"`
}
