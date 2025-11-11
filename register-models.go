package main

import (
	"avileads-web/models"
	"github.com/astaxie/beego/orm"
)

func initOrmModels() {
	orm.RegisterModel(new(models.Roles))
	orm.RegisterModel(new(models.OpenAIMessages))
	orm.RegisterModel(new(models.Region))
	orm.RegisterModel(new(models.CityDistrict))
	orm.RegisterModel(new(models.District))
	orm.RegisterModel(new(models.CityType))
	orm.RegisterModel(new(models.BuildingType))
	orm.RegisterModel(new(models.StatusDict))
	orm.RegisterModel(new(models.Address))
	orm.RegisterModel(new(models.AddressBooking))
	orm.RegisterModel(new(models.Booking))
	orm.RegisterModel(new(models.BookingCity))
	orm.RegisterModel(new(models.BookingDistrict))
	orm.RegisterModel(new(models.Metro))
	orm.RegisterModel(new(models.BookingMetro))
	orm.RegisterModel(new(models.RkRekl))
	orm.RegisterModel(new(models.RkBuyingType))
	orm.RegisterModel(new(models.RkProfession))
	orm.RegisterModel(new(models.RkProvider))
	orm.RegisterModel(new(models.RkTrafficSource))
	orm.RegisterModel(new(models.ExportType))
	orm.RegisterModel(new(models.AuditLog))
	orm.RegisterModel(new(models.AuditModel))
	orm.RegisterModel(new(models.AuditAction))
	orm.RegisterModel(new(models.Profession))
	orm.RegisterModel(new(models.ProfessionConflict))
}
