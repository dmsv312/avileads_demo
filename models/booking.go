package models

import (
	"avileads-web/utils"
	"github.com/astaxie/beego/orm"
	"time"
)

type Region struct {
	Id   int    `orm:"column(id);auto"`
	Name string `orm:"column(name)"`
}

type District struct {
	Id     int     `orm:"column(id);auto"`
	Region *Region `orm:"column(region_id);rel(fk)"`
	Name   string  `orm:"column(name)"`
}

type City struct {
	Id       int       `orm:"column(id);auto;pk"`
	Name     string    `orm:"column(name)"`
	Region   *Region   `orm:"column(region_id);null;rel(fk)"`
	District *District `orm:"column(district_id);null;rel(fk)"`
	Type     *CityType `orm:"column(type_id);null;rel(fk)"`
}

type CityDistrict struct {
	Id   int    `form:"-"`
	Name string `orm:"column(name)"`
	City *City  `orm:"column(city_id);null;rel(fk)"`
}

type Metro struct {
	Id           int           `orm:"column(id);pk;auto"`
	OsmId        int64         `orm:"column(osm_id)"`
	Name         string        `orm:"column(name)"`
	Location     string        `orm:"column(location);type(geometry);"`
	CityDistrict *CityDistrict `orm:"column(city_district_id);rel(fk)"`
}

type CityType struct {
	Id   int    `form:"-"`
	Name string `orm:"column(name)"`
}

type BuildingType struct {
	Id   int    `orm:"column(id);auto"`
	Name string `orm:"column(name)"`
}

type StatusDict struct {
	Id   int    `orm:"column(id);pk"`
	Name string `orm:"column(name)"`
}

type Address struct {
	Id           int           `orm:"column(id);auto"`
	OsmId        int           `orm:"column(osm_id);null"`
	OriginalId   int           `orm:"column(original_id);"`
	Name         string        `orm:"column(name);null"`
	Location     string        `orm:"column(location);type(geometry)"`
	BuildingType *BuildingType `orm:"column(building_type_id);rel(fk)"`
	PostalCode   string        `orm:"column(postal_code);null"`
	City         *City         `orm:"column(city_id);rel(fk)"`
	CityDistrict *CityDistrict `orm:"column(city_district_id);rel(fk);null"`
	Street       string        `orm:"column(street);null"`
	House        string        `orm:"column(house);null"`
	FullAddress  string        `orm:"column(full_address);null"`
}

type AddressBooking struct {
	Id          int         `orm:"column(id);auto"`
	Address     *Address    `orm:"column(address_id);rel(fk)"`
	Booking     *Booking    `orm:"column(booking_id);rel(fk)"`
	IsActive    bool        `orm:"column(is_active)"`
	DaysNumber  int         `orm:"column(days_number)"`
	StartDate   time.Time   `orm:"column(start_date);type(date)"`
	EndDate     time.Time   `orm:"column(end_date);type(date)"`
	PostingDate time.Time   `orm:"column(posting_date);type(timestamp)"`
	Comment     string      `orm:"column(comment);null"`
	AccountId   int         `orm:"column(account_id);null"`
	Profession  *Profession `orm:"column(profession_id);null;rel(fk)"`
}

type Booking struct {
	Id           int       `orm:"column(id);auto"`
	VacancyId    int       `orm:"column(vacancy_id)"`
	BookingCount int       `orm:"column(booking_count)"`
	DaysNumber   int       `orm:"column(days_number)"`
	IsActive     bool      `orm:"column(is_active)"`
	StartDate    time.Time `orm:"column(start_date);type(timestamp)"`
	EndDate      time.Time `orm:"column(end_date);type(timestamp)"`
	PostingDate  time.Time `orm:"column(posting_date);type(timestamp)"`
	User         *Users    `orm:"column(user_id);rel(fk)"`
	Status       string    `orm:"column(status)"`
	ResultLink   string    `orm:"column(result_link)"`
	GroupUUID    string    `orm:"column(group_uuid)"`

	AddressBookings  []*AddressBooking  `orm:"reverse(many)"`
	BookingCities    []*BookingCity     `orm:"reverse(many)"`
	BookingDistricts []*BookingDistrict `orm:"reverse(many)"`

	AccountId  int         `orm:"column(account_id);null"`
	Profession *Profession `orm:"column(profession_id);null;rel(fk)"`
}

type BookingCity struct {
	Id      int      `orm:"column(id);auto"`
	Booking *Booking `orm:"column(booking_id);rel(fk)"`
	City    *City    `orm:"column(city_id);rel(fk)"`
}

type BookingDistrict struct {
	Id       int           `orm:"column(id);auto"`
	Booking  *Booking      `orm:"column(booking_id);rel(fk)"`
	District *CityDistrict `orm:"column(district_id);rel(fk)"`
}

type BookingMetro struct {
	Id      int      `orm:"column(id);pk;auto"`
	Booking *Booking `orm:"column(booking_id);rel(fk)"`
	Metro   *Metro   `orm:"column(metro_id);rel(fk)"`
}

type Profession struct {
	Id   int    `orm:"column(id);pk;auto"`
	Code string `orm:"column(code);unique"`
	Name string `orm:"column(name)"`
}

func (p *Profession) TableName() string { return "geo_profession" }

type ProfessionConflict struct {
	Id            int         `orm:"column(id);pk;auto"`
	Profession    *Profession `orm:"column(profession_id);rel(fk)"`
	ConflictsWith *Profession `orm:"column(conflicts_with_id);rel(fk)"`
}

func (pc *ProfessionConflict) TableName() string { return "geo_profession_conflict" }

func (r *Region) TableName() string {
	return "geo_region"
}

func (s *District) TableName() string {
	return "geo_district"
}

func (c *City) TableName() string {
	return "geo_city"
}

func (c *CityDistrict) TableName() string {
	return "geo_city_district"
}

func (c *CityType) TableName() string {
	return "geo_city_type"
}

func (bt *BuildingType) TableName() string {
	return "geo_building_type"
}

func (s *StatusDict) TableName() string { return "geo_status_dict" }

func (a *Address) TableName() string {
	return "geo_address"
}

func (ab *AddressBooking) TableName() string {
	return "geo_address_booking"
}

func (b *Booking) TableName() string {
	return "geo_booking"
}

func (bc *BookingCity) TableName() string {
	return "geo_booking_city"
}

func (bd *BookingDistrict) TableName() string {
	return "geo_booking_district"
}

func (m *Metro) TableName() string {
	return "geo_metro"
}

func (b *BookingMetro) TableName() string {
	return "geo_booking_metro"
}

func GetAddressBooking(addressBookingId int) (*AddressBooking, error) {
	o := utils.CreateDefaultDbContext()
	addressBooking := &AddressBooking{Id: addressBookingId}
	err := o.Read(addressBooking)
	return addressBooking, err
}

func UpdateAddressBooking(addressBooking *AddressBooking) error {
	o := orm.NewOrm()
	_, err := o.Update(addressBooking)
	return err
}

func GetBooking(bookingId int) (*Booking, error) {
	o := utils.CreateDefaultDbContext()
	booking := &Booking{Id: bookingId}
	err := o.Read(booking)
	return booking, err
}

func UpdateBooking(booking *Booking) error {
	o := orm.NewOrm()
	_, err := o.Update(booking)
	return err
}

func CreateBookingCity(bookingCity *BookingCity) error {
	_, err := utils.CreateDefaultDbContext().Insert(bookingCity)
	return err
}

func GetBookingCity(bookingCityId int) (*BookingCity, error) {
	o := utils.CreateDefaultDbContext()
	bookingCity := &BookingCity{Booking: &Booking{Id: bookingCityId}}
	err := o.Read(bookingCity, "booking_id")
	return bookingCity, err
}

func UpdateBookingCity(bookingCity *BookingCity) error {
	o := orm.NewOrm()
	_, err := o.Update(bookingCity)
	return err
}

func DeleteBookingCity(bookingCityId int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&BookingCity{Booking: &Booking{Id: bookingCityId}})
	return err
}

func CreateBookingDistrict(bookingDistrict *BookingDistrict) error {
	_, err := utils.CreateDefaultDbContext().Insert(bookingDistrict)
	return err
}

func GetBookingDistrict(bookingDistrictId int) (*BookingDistrict, error) {
	o := utils.CreateDefaultDbContext()
	bookingDistrict := &BookingDistrict{Booking: &Booking{Id: bookingDistrictId}}
	err := o.Read(bookingDistrict, "booking_id")
	return bookingDistrict, err
}

func UpdateBookingDistrict(bookingDistrict *BookingDistrict) error {
	o := orm.NewOrm()
	_, err := o.Update(bookingDistrict)
	return err
}

func DeleteBookingDistrict(bookingDistrictId int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&BookingDistrict{Booking: &Booking{Id: bookingDistrictId}})
	return err
}

func UpdateAddressBookingsByBooking(bookingId int, booking *Booking) error {
	o := orm.NewOrm()
	_, err := o.QueryTable(new(AddressBooking)).
		Filter("Booking__Id", bookingId).
		Update(orm.Params{
			"EndDate":     booking.EndDate,
			"PostingDate": booking.PostingDate,
			"DaysNumber":  booking.DaysNumber,
		})
	return err
}
