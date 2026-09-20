package models

type Region struct {
	Name string `gorm:"column:name;type:varchar(200)" json:"name"`
	Code string `gorm:"column:code;type:varchar(20)" json:"code"`

	MainBaseModel
}

type RegionHistory struct {
	Region
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
