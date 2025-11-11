package controllers

import (
	"avileads-web/application"
	"avileads-web/models"
	"avileads-web/services"
	"avileads-web/transport"
	"avileads-web/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type VacancyAdminController struct {
	beego.Controller
}

const pageSize = 20

func (c *VacancyAdminController) VacancyAdminSingleView() {
	c.TplName = "dictionary/vacancy_admin/vacancy_admin.html"
	ClientId := (&utils.HttpContext{Context: c.Ctx}).GetAuthClientId()
	vacancyId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	//flash := beego.ReadFromRequest(&c.Controller)

	o := orm.NewOrm()
	o.Using("default")
	vacancy := models.Vacancy_avito{Id: vacancyId, ClientId: ClientId}
	var quest []*models.Questionnaire
	num, err := o.QueryTable("questionnaire").Filter("clientId", ClientId).OrderBy("name").All(&quest)
	if err != orm.ErrNoRows && num > 0 {
		c.Data["quest_types"] = quest
	}
	if o.Read(&vacancy) == nil {
		_, _ = o.LoadRelated(&vacancy, "quest_id")
		c.Data["records"] = vacancy
		if vacancy.CountOfMessagesAfterFinishedDialog.Valid {
			c.Data["countOfMessagesAfterFinishedDialog"] = vacancy.CountOfMessagesAfterFinishedDialog.Int64
		} else {
			c.Data["countOfMessagesAfterFinishedDialog"] = ""
		}
	}

	var Vacancy_2_Offer []*models.Vacancy_2_Offer
	num, err = o.QueryTable("vacancy_2_vacancy").Filter("VacancyAvito", vacancyId).OrderBy("Id").RelatedSel().All(&Vacancy_2_Offer)
	if err != orm.ErrNoRows && num > 0 {
		for _, v2o := range Vacancy_2_Offer {
			filters, err := models.GetOfferSQLFilters(o, v2o.Offer.Id, vacancy.ClientId)
			if err != nil {
				continue
			}
			v2o.Offer.OfferFilters = filters

			sqlFilters, err := models.GetSQLFiltersForOffer(o, v2o.Offer.Id, vacancy.ClientId)
			if err != nil {
				continue
			}
			v2o.Offer.SqlFilters = sqlFilters

			clientSettings, err := models.GetClientOfferSettings(o, v2o.Offer.Id, vacancy.ClientId)
			if err != nil {
				continue
			}
			v2o.Offer.ClientSetting = clientSettings
		}

		c.Data["offers"] = Vacancy_2_Offer
	}

	vacancyAvitoOpenAI, err := models.GetVacancyAvitoOpenAIById(vacancyId)
	if err != orm.ErrNoRows {
		c.Data["vacancy_avito_openai"] = vacancyAvitoOpenAI
	}

	var vacancyText string
	err = o.Raw("SELECT vacancy_text FROM vacancy_avito_view WHERE id = ?", vacancyId).QueryRow(&vacancyText)
	if err != orm.ErrNoRows && err == nil {
		c.Data["information"] = vacancyText
	} else {
		c.Data["information"] = "Информация по вакансии не найдена."
	}
}

func (c *VacancyAdminController) VacancyAdminUpdate() {
	o := orm.NewOrm()
	o.Using("default")
	flash := beego.NewFlash()
	referer := c.Ctx.Request.Header.Get("Referer")
	index := strings.Index(referer, "/dict")

	vacancyIdStr := c.GetString("vacancyId")
	vacancyId, err := strconv.Atoi(vacancyIdStr)
	if err != nil {
		c.Data["errors"] = "Ошибка получения id вакансии"
		c.Redirect(referer[index:], 302)
		return
	}

	vacancy := models.Vacancy_avito{Id: vacancyId}
	if o.Read(&vacancy) != nil {
		c.Data["errors"] = fmt.Sprintf("Вакансия не найдена")
		c.Redirect(referer[index:], 302)
		return
	}

	vacancy.LoadToZp = c.GetString("LoadToZP") == "on"

	questionnaireId, err := strconv.Atoi(c.GetString("questionnaire"))
	if err == nil {
		vacancy.Quest = &models.Questionnaire{Id: questionnaireId}
	}

	_, err = o.Update(&vacancy)
	if err != nil {
		c.Data["errors"] = fmt.Sprintf("Не удалось обновить вакансию: %s", err.Error())
		c.Redirect(referer[index:], 302)
		return
	}

	flash.Notice("Вакансия успешно обновлена")
	flash.Store(&c.Controller)
	c.Redirect(referer[index:], 302)
}

func (c *VacancyAdminController) VacancyCopy() {
	o := orm.NewOrm()
	o.Using("default")

	clientId := (&utils.HttpContext{Context: c.Ctx}).GetAuthClientId()
	vacancyID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	oldVacancy := models.Vacancy_avito{Id: vacancyID, ClientId: clientId}
	if err := o.Read(&oldVacancy); err == nil {
		copyData := &models.Vacancy_avito{
			LoadToZp:                           oldVacancy.LoadToZp,
			HelloText:                          oldVacancy.HelloText,
			ByeText:                            oldVacancy.ByeText,
			PostProcessingMessages:             oldVacancy.PostProcessingMessages,
			DialogLifeTimeInMinutes:            oldVacancy.DialogLifeTimeInMinutes,
			FollowUpMessageIntervalInMinutes:   oldVacancy.FollowUpMessageIntervalInMinutes,
			CountOfMessagesAfterFinishedDialog: oldVacancy.CountOfMessagesAfterFinishedDialog,
			Name:                               oldVacancy.Name,
			OpenAI_Support:                     oldVacancy.OpenAI_Support,
		}
		c.SetSession("vacancyCopyData", copyData)
	}

	c.Redirect("/dict/vacancy/edit/0?copy=1", 302)
}

func (c *VacancyAdminController) VacancyAdminView() {
	c.TplName = "dictionary/vacancy_admin/vacancy_list.html"
	ClientId := (&utils.HttpContext{Context: c.Ctx}).GetAuthClientId()
	flash := beego.ReadFromRequest(&c.Controller)

	if ok := flash.Data["error"]; ok != "" {
		c.Data["errors"] = ok
	}

	if ok := flash.Data["notice"]; ok != "" {
		c.Data["notices"] = ok
	}
	o := orm.NewOrm()
	o.Using("default")

	var vacancyOpenAI []orm.Params
	_, err := o.Raw("Select * from vacancy_openai_view where \"clientId\"=?", ClientId).Values(&vacancyOpenAI)
	if err != orm.ErrNoRows {
		var openAIFull, openAIIncomplete []orm.Params
		for _, vac := range vacancyOpenAI {
			if vac["offers"] != nil {
				openAIFull = append(openAIFull, vac)
			} else {
				openAIIncomplete = append(openAIIncomplete, vac)
			}
		}

		c.Data["open_ai_full"] = openAIFull
		c.Data["open_ai_incomplete"] = openAIIncomplete
	}

	var vacancy []orm.Params
	_, err = o.Raw("Select * from vacancy_view where \"clientId\"=?", ClientId).Values(&vacancy)
	var fullVacancies, incompleteVacancies []orm.Params
	if err != orm.ErrNoRows {
		for _, vac := range vacancy {
			if vac["questname"] != nil && vac["offers"] != nil {
				fullVacancies = append(fullVacancies, vac)
			} else {
				incompleteVacancies = append(incompleteVacancies, vac)
			}
		}
	}

	if fullData, ok := c.Data["open_ai_full"].([]orm.Params); ok {
		fullVacancies = append(fullVacancies, fullData...)
	}

	if incompleteData, ok := c.Data["open_ai_incomplete"].([]orm.Params); ok {
		incompleteVacancies = append(incompleteVacancies, incompleteData...)
	}

	c.Data["records_full"] = fullVacancies
	c.Data["records_incomplete"] = incompleteVacancies
}

func (c *VacancyAdminController) VacancyListAdminView() {
	flash := beego.ReadFromRequest(&c.Controller)

	if ok := flash.Data["error"]; ok != "" {
		c.Data["errors"] = ok
	}

	if ok := flash.Data["notice"]; ok != "" {
		c.Data["notices"] = ok
	}

	page, _ := c.GetInt("page", 1)
	if page < 1 {
		page = 1
	}
	listType := c.GetString("list", "vacancy")
	hideArchived, _ := c.GetInt("hide_archived", 1)
	searchString := strings.TrimSpace(c.GetString("search", ""))
	avitoSearch := strings.TrimSpace(c.GetString("avitoIdSearch", ""))

	var avitoIDs []string
	if avitoSearch != "" {
		for _, s := range strings.Split(avitoSearch, ",") {
			if id := strings.TrimSpace(s); id != "" {
				avitoIDs = append(avitoIDs, id)
			}
		}
	}

	ClientId := (&utils.HttpContext{Context: c.Ctx}).GetAuthClientId()
	o := orm.NewOrm()
	o.Using("default")

	buildWhere := func() (join string, where string, args []interface{}) {
		where = ` WHERE v."clientId" = ?`
		args = []interface{}{ClientId}

		if hideArchived == 1 {
			where += ` AND v.load_to_zp = true`
		}
		if len(searchString) >= 3 {
			where += `
			  AND (
				   lower(v.name) LIKE lower(?)
				OR cast(v.id AS text) ILIKE ?
			  )`
			like := "%" + searchString + "%"
			args = append(args, like, like)
		}

		if len(avitoIDs) > 0 {
			ph := strings.Repeat("?,", len(avitoIDs))
			ph = ph[:len(ph)-1]

			join = ` JOIN avito_vacancy_ids avi ON avi.vacancy_avito_id = v.id`
			where += ` AND avi.avito_id IN (` + ph + `)`
			for _, id := range avitoIDs {
				args = append(args, id)
			}
		}
		return
	}

	fetch := func(view string, needTotal bool) (
		rows []orm.Params,
		pages int64,
		err error) {

		join, where, args := buildWhere()

		baseSQL := ` FROM ` + view + ` v` + join + where

		if needTotal {
			var total int64
			if err = o.Raw(`SELECT COUNT(DISTINCT v.id)`+baseSQL, args...).
				QueryRow(&total); err != nil {
				return
			}

			offset := (page - 1) * pageSize
			argsData := append(args, pageSize, offset)

			_, err = o.Raw(`
            SELECT DISTINCT v.*
            `+baseSQL+`
            ORDER BY v.id DESC
            LIMIT ? OFFSET ?`, argsData...).Values(&rows)

			pages = (total + int64(pageSize) - 1) / int64(pageSize)

		} else {
			argsData := append(args, pageSize)
			_, err = o.Raw(`
            SELECT DISTINCT v.*
            `+baseSQL+`
            ORDER BY v.id DESC
            LIMIT ?`, argsData...).Values(&rows)
		}
		return
	}

	var (
		vacRows, aiRows   []orm.Params
		pagesVac, pagesAI int64
	)

	if listType == "vacancy" {
		vacRows, pagesVac, _ = fetch("rk_view", true)
		aiRows, _, _ = fetch("rk_openai_view", false)
	} else {
		aiRows, pagesAI, _ = fetch("rk_openai_view", true)
		vacRows, _, _ = fetch("rk_view", false)
	}

	const maxButtons = 7
	if listType == "vacancy" {
		c.Data["page_list"] = utils.BuildPageList(int(pagesVac), page, maxButtons)
	} else {
		c.Data["page_list_openai"] = utils.BuildPageList(int(pagesAI), page, maxButtons)
	}

	c.Data["records"] = vacRows
	c.Data["records_open_ai"] = aiRows
	c.Data["pages_total"] = pagesVac
	c.Data["pages_open_ai_total"] = pagesAI
	c.Data["active_tab"] = listType
	c.Data["hide_archived"] = hideArchived
	c.Data["page"] = page
	c.Data["search"] = searchString
	c.Data["avitoIdSearch"] = avitoSearch

	c.TplName = "dictionary/vacancy_admin/vacancy_list.html"

}

func (c *VacancyAdminController) VacancyAdminUpdateNew() {
	var request transport.PromoCampaignPayload

	_ = c.Ctx.Request.ParseMultipartForm(100 << 20)

	payload := c.Ctx.Request.FormValue("payload")
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		c.CustomAbort(400, "invalid payload: "+err.Error())
		return
	}

	file, hdr, _ := c.GetFile("avitoIdsFile")
	switch {
	case file != nil:
		avitoIds, err := services.ExtractAvitoIds(file, hdr)
		if err != nil {
			c.CustomAbort(400, "invalid payload: "+err.Error())
			return
		}
		request.Vacancy.IdAvitoList = avitoIds

	case strings.TrimSpace(request.ManualAvitoIds) != "":
		request.Vacancy.IdAvitoList = utils.CleanAvitoIds(request.ManualAvitoIds)

	default:

	}

	clientId := (&utils.HttpContext{Context: c.Ctx}).GetAuthClientId()
	userId := (&utils.HttpContext{Context: c.Ctx}).GetAuthUserId()

	promoCampaign := application.NewPromoCampaign(orm.NewOrm(), clientId, userId)
	vacancyId, err := promoCampaign.Process(request)
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{
		"status":    "ok",
		"vacancyId": vacancyId,
	}
	c.ServeJSON()
}

func (c *VacancyAdminController) VacancyAdminViewNew() {
	c.TplName = "dictionary/vacancy_admin/vacancy_admin.html"

	vacancyId, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	copyFrom, _ := c.GetInt("copyFrom")

	sourceId := vacancyId
	if vacancyId == 0 && copyFrom > 0 {
		sourceId = copyFrom
	}

	userId := (&utils.HttpContext{Context: c.Ctx}).GetAuthUserId()
	clientId := (&utils.HttpContext{Context: c.Ctx}).GetAuthClientId()

	o := orm.NewOrm()
	promoCampaign := application.NewPromoCampaign(o, clientId, userId)
	response, err := promoCampaign.Get(sourceId)
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	if vacancyId == 0 && copyFrom > 0 {
		response.Vacancy.Id = 0
		for i := range response.Offers {
			response.Offers[i].Id = 0
		}
		response.Questionnaire.Id = 0
	}

	var vacancyText string
	err = o.Raw("SELECT vacancy_text FROM vacancy_avito_view WHERE id = ?", sourceId).QueryRow(&vacancyText)
	if err != orm.ErrNoRows && err == nil {
		c.Data["information"] = vacancyText
	} else {
		c.Data["information"] = "Информация по вакансии не найдена."
	}

	var total int64
	err = o.Raw(`
        SELECT COUNT(*) 
        FROM avito_vacancy_ids 
        WHERE vacancy_avito_id = ?`,
		vacancyId,
	).QueryRow(&total)

	if err != orm.ErrNoRows {
		c.Data["avito_ids"] = fmt.Sprintf("К РК привязано %d ID Авито", total)
	}

	c.Data["records"] = response.Vacancy
	c.Data["offers"] = response.Offers
	c.Data["quest"] = response.Questionnaire
	c.Data["questions"] = response.Questionnaire.Questions
}

func (c *VacancyAdminController) CheckAvitoIds() {
	if err := c.Ctx.Request.ParseMultipartForm(100 << 20); err != nil {
		c.CustomAbort(400, "cant parse multipart: "+err.Error())
		return
	}

	var form transport.CheckDocumentPayload
	if err := c.ParseForm(&form); err != nil {
		c.CustomAbort(400, "cant parse form: "+err.Error())
		return
	}

	vacancyId, _ := c.GetInt("vacancyId", 0)

	file, hdr, _ := c.GetFile("avitoIdsFile")

	newCnt, err := services.PreviewNewAvitoIds(
		vacancyId,
		form.InputType,
		file,
		hdr,
		form.ManualAvitoIds,
		orm.NewOrm(),
	)

	if err != nil {
		c.CustomAbort(400, err.Error())
		return
	}

	c.Data["json"] = map[string]int{"new": newCnt}
	c.ServeJSON()
}
