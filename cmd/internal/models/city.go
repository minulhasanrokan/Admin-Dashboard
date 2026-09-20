package models

type City struct {
	Name        string `gorm:"column:name;type:varchar(200)" json:"name"`
	Code        string `gorm:"column:code;type:varchar(20)" json:"code"`
	StateCode   string `gorm:"column:state_code;type:varchar(20)" json:"state_code"`
	CountryCode string `gorm:"column:country_code;type:varchar(20)" json:"country_code"`
	MainBaseModel
}

type CityHistory struct {
	City
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
