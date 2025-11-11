package httphandler

import (
	"avileads-web/models"
	"avileads-web/services"
	"avileads-web/services/rk"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"io"
	"net/http"
	"os"
	"strconv"
)

func GetAllSqlFilters(w http.ResponseWriter, r *http.Request) {
	o := orm.NewOrm()
	o.Using("default")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sqlFilters, _ := models.GetAllSqlFilters(o)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sqlFilters)
}

func GetClients(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	clients, _ := models.GetClients()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

func GetExportTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	exportTypes, _ := models.GetExportTypes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exportTypes)
}

func GetQuestionTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	questionTypes, _ := models.GetQuestionTypes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questionTypes)
}

func GetVacancyText() {
	//var vacancyText string
	//err = o.Raw("SELECT vacancy_text FROM vacancy_avito_view WHERE id = ?", vacancyId).QueryRow(&vacancyText)
	//if err != orm.ErrNoRows && err == nil {
	//	c.Data["information"] = vacancyText
	//} else {
	//	c.Data["information"] = "Информация по вакансии не найдена."
	//}
}

func GetBuyingTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	service := rk.NewDictAggregateService(orm.NewOrm(), nil)
	buyingTypes, err := service.GetBuyingTypes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buyingTypes)
}
func GetRekls(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	service := rk.NewDictAggregateService(orm.NewOrm(), nil)
	rekls, err := service.GetRekls()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rekls)
}

func GetTrafficSources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	service := rk.NewDictAggregateService(orm.NewOrm(), nil)
	trafficSources, err := service.GetTrafficSources()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trafficSources)
}

func GetProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	service := rk.NewDictAggregateService(orm.NewOrm(), nil)
	providers, err := service.GetProviders()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

func GetProfessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	svc := rk.NewDictAggregateService(orm.NewOrm(), nil)
	professions, err := svc.GetProfessions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(professions)
}

func GetAvitoIdsExcel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vacancyIdStr := r.URL.Query().Get("vacancyId")
	if vacancyIdStr == "" {
		http.Error(w, "vacancyId обязателен", http.StatusBadRequest)
		return
	}
	vacancyId, err := strconv.Atoi(vacancyIdStr)
	if err != nil || vacancyId <= 0 {
		http.Error(w, "vacancyId должен быть положительным числом", http.StatusBadRequest)
		return
	}

	o := orm.NewOrm()
	path, err := services.ExportAvitoIdsToExcel(vacancyId, o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(path)

	f, err := os.Open(path)
	if err != nil {
		http.Error(w, "не удалось открыть файл", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="avito_ids_%d.xlsx"`, vacancyId))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}
