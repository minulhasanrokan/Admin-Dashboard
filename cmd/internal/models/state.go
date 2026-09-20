package models

type State struct {
	Name        string `gorm:"column:name;type:varchar(200)" json:"name"`
	Code        string `gorm:"column:code;type:varchar(20)" json:"code"`
	CountryCode string `gorm:"column:country_code;type:varchar(20)" json:"country_code"`
	MainBaseModel
}

type StateHistory struct {
	State
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
