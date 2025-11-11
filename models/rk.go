package models

type RkBuyingType struct {
	Id        int    `orm:"column(id);pk;auto"`
	Name      string `orm:"column(name)"`
	ShortName string `orm:"column(short_name)"`
}

type RkProfession struct {
	Id        int    `orm:"column(id);pk;auto"`
	Name      string `orm:"column(name)"`
	ShortName string `orm:"column(short_name)"`
	BasicName string `orm:"column(basic_name)"`
}

type RkProvider struct {
	Id        int    `orm:"column(id);pk;auto"`
	Name      string `orm:"column(name)"`
	ShortName string `orm:"column(short_name)"`
	Type      string `orm:"column(type)"`
}

type RkRekl struct {
	Id        int    `orm:"column(id);pk;auto"`
	Name      string `orm:"column(name)"`
	ShortName string `orm:"column(short_name)"`
}

type RkTrafficSource struct {
	Id        int    `orm:"column(id);pk;auto"`
	Name      string `orm:"column(name)"`
	ShortName string `orm:"column(short_name)"`
}

func (RkTrafficSource) TableName() string { return "rk_dict_traffic_source" }

func (RkProfession) TableName() string { return "rk_dict_prof" }

func (RkRekl) TableName() string { return "rk_dict_rekl" }

func (RkProvider) TableName() string { return "rk_dict_provider" }

func (RkBuyingType) TableName() string { return "rk_dict_buying_type" }
