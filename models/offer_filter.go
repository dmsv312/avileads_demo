package models

type OfferFilter struct {
	Id          int `orm:"column(id);pk;auto"`
	OfferId     int `orm:"column(offer_id)"`
	SqlFilterId int `orm:"column(sql_filter_id)"`
	ClientId    int `orm:"column(client_id)"`
}

func (OfferFilter) TableName() string { return "offer_filters" }
