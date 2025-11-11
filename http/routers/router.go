package routers

import (
	"avileads-web/http/controllers"
	"avileads-web/http/httphandler"
	"avileads-web/models"
	"avileads-web/utils"
	"net/http"
	"regexp"

	"github.com/astaxie/beego"
)

type RouteAccess struct {
	Route          string
	AvailableRules []string
}

var routeAccesses = []RouteAccess{
	{"/", []string{}},
	{"/login", []string{}},
	{"/logout", []string{}},

	{"/data/leads", []string{models.Data_RuleKey}},
	{"/analytics/state", []string{models.Analytics_RuleKey}},
	{"/admin/users", []string{models.Admin_Users_RuleKey}},
	{"/admin/roles", []string{models.Admin_Roles_RuleKey}},

	{"/dict/vacancy/copy/:id([0-9]+)", []string{models.Dict_Vacancy_RuleKey}},
	{"/dict/vacancy_admin/edit/([0-9]+)", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/dict/vacancy_admin/update", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/dict/vacancy_admin", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/dict/vacancy_admin/save_all", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/rk/check_avito_ids", []string{models.Dict_Vacancy_Admin_RuleKey}},

	{"/rk/filters", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/rk/clients", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/rk/export_types", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/rk/question_types", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/rk/buying_types", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/rk/professions", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/rk/providers", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/rk/rekls", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/rk/traffic_sources", []string{models.Dict_Vacancy_Admin_RuleKey}},
	{"/rk/avito_ids", []string{models.Dict_Vacancy_Admin_RuleKey}},

	{"/booking/create", []string{models.Randomizer_RuleKey}},
	{"/booking/manage", []string{models.Randomizer_RuleKey}},
	{"/booking/directory", []string{models.Randomizer_RuleKey}},
	{"/booking/directory/cities", []string{models.Randomizer_RuleKey}},
	{"/booking/directory/districts", []string{models.Randomizer_RuleKey}},
	{"/booking", []string{models.Randomizer_RuleKey}},
	{"/booking/regions", []string{models.Randomizer_RuleKey}},
	{"/booking/cities", []string{models.Randomizer_RuleKey}},
	{"/booking/city_districts", []string{models.Randomizer_RuleKey}},
	{"/booking/metro", []string{models.Randomizer_RuleKey}},
	{"/booking/vacancies", []string{models.Randomizer_RuleKey}},
	{"/booking/professions", []string{models.Randomizer_RuleKey}},
	{"/booking/accounts", []string{models.Randomizer_RuleKey}},
	{"/booking/save_bookings", []string{models.Randomizer_RuleKey}},
	{"/booking/list", []string{models.Randomizer_RuleKey}},
	{"/booking/addresses", []string{models.Randomizer_RuleKey}},
	{"/booking/delete_addresses", []string{models.Randomizer_RuleKey}},
	{"/booking/delete_bookings", []string{models.Randomizer_RuleKey}},
	{"/booking/update_booking", []string{models.Randomizer_RuleKey}},
	{"/booking/update_address", []string{models.Randomizer_RuleKey}},
	{"/booking/job", []string{models.Randomizer_RuleKey}},
	{"/booking/export", []string{models.Randomizer_RuleKey}},
}

func AllRouteAccess() []RouteAccess {
	copiedRules := make([]RouteAccess, len(routeAccesses))
	copy(copiedRules, routeAccesses)
	return copiedRules
}

func HasAccessToPath(path string, currentRules []string) bool {
	var routeAccess RouteAccess
	hasExistRoute := false

	for _, x := range routeAccesses {
		matched, err := regexp.MatchString("^"+x.Route+"$", path)
		if err != nil {
			continue
		}
		if matched {
			routeAccess = x
			hasExistRoute = true
			break
		}
	}

	if !hasExistRoute {
		return false
	}

	// For all access
	if len(routeAccess.AvailableRules) == 0 {
		return true
	}

	ruleExists := func(rule string, availableRules []string) bool {
		for _, r := range availableRules {
			if r == rule {
				return true
			}
		}
		return false
	}

	for _, rule := range currentRules {
		if ruleExists(rule, routeAccess.AvailableRules) {
			return true
		}
	}

	return false
}

var sessionManager = utils.NewSessionManager()

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/vacancy_avito/delete/:id([0-9]+)", &controllers.DictionaryController{}, "*:Delete")
	beego.Router("/login", &controllers.MainController{}, "get:LoginGet")
	beego.Router("/login", &controllers.MainController{}, "post:LoginPost")
	beego.Router("/logout", &controllers.MainController{}, "get:LogoutPost")

	beego.Router("/dict/vacancy/copy/:id([0-9]+)", &controllers.VacancyAdminController{}, "get:VacancyCopy")
	beego.Router("/dict/vacancy_admin/edit/:id([0-9]+)", &controllers.VacancyAdminController{}, "get:VacancyAdminViewNew")
	beego.Router("/dict/vacancy_admin/update", &controllers.VacancyAdminController{}, "post:VacancyAdminUpdate")
	beego.Router("/dict/vacancy_admin", &controllers.VacancyAdminController{}, "get:VacancyListAdminView")
	beego.Router("/dict/vacancy_admin/save_all", &controllers.VacancyAdminController{}, "post:VacancyAdminUpdateNew")
	beego.Router("/rk/check_avito_ids", &controllers.VacancyAdminController{}, "post:CheckAvitoIds")

	beego.Handler("/rk/filters", http.HandlerFunc(httphandler.GetAllSqlFilters))
	beego.Handler("/rk/clients", http.HandlerFunc(httphandler.GetClients))
	beego.Handler("/rk/export_types", http.HandlerFunc(httphandler.GetExportTypes))
	beego.Handler("/rk/question_types", http.HandlerFunc(httphandler.GetQuestionTypes))
	beego.Handler("/rk/buying_types", http.HandlerFunc(httphandler.GetBuyingTypes))
	beego.Handler("/rk/professions", http.HandlerFunc(httphandler.GetProfessions))
	beego.Handler("/rk/providers", http.HandlerFunc(httphandler.GetProviders))
	beego.Handler("/rk/rekls", http.HandlerFunc(httphandler.GetRekls))
	beego.Handler("/rk/traffic_sources", http.HandlerFunc(httphandler.GetTrafficSources))
	beego.Handler("/rk/avito_ids", http.HandlerFunc(httphandler.GetAvitoIdsExcel))

	beego.Router("/booking/create", &controllers.BookingController{}, "get:Home")
	beego.Router("/booking/manage", &controllers.BookingController{}, "get:List")
	beego.Router("/booking/directory", &controllers.BookingController{}, "get:Directory")
	beego.Handler("/booking/directory/cities", http.HandlerFunc(httphandler.GetDirectoryCities))
	beego.Handler("/booking/directory/districts", http.HandlerFunc(httphandler.GetDirectoryDistricts))
	beego.Handler("/booking/regions", http.HandlerFunc(httphandler.GetRegions))
	beego.Handler("/booking/cities", http.HandlerFunc(httphandler.GetCities))
	beego.Handler("/booking/city_districts", http.HandlerFunc(httphandler.GetCityDistricts))
	beego.Handler("/booking/metro", http.HandlerFunc(httphandler.GetMetro))
	beego.Handler("/booking/vacancies", http.HandlerFunc(httphandler.GetVacancies))
	beego.Handler("/booking/professions", http.HandlerFunc(httphandler.GetProfessionsForBooking))
	beego.Handler("/booking/accounts", http.HandlerFunc(httphandler.GetAvitoAccountsForBooking))
	beego.Handler("/booking/save_bookings", http.HandlerFunc(httphandler.SaveBookings))
	beego.Handler("/booking/list", http.HandlerFunc(httphandler.GetBookings))
	beego.Handler("/booking/addresses", http.HandlerFunc(httphandler.GetBookingAddresses))
	beego.Handler("/booking/delete_bookings", http.HandlerFunc(httphandler.DeleteBookings))
	beego.Handler("/booking/delete_addresses", http.HandlerFunc(httphandler.DeleteBookingAddresses))
	beego.Handler("/booking/update_booking", http.HandlerFunc(httphandler.UpdateBooking))
	beego.Handler("/booking/update_address", http.HandlerFunc(httphandler.UpdateAddressBooking))
	beego.Handler("/booking/export", http.HandlerFunc(httphandler.ExportBookings))
}
