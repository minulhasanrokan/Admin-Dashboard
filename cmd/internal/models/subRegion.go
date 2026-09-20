package models

type SubRegion struct {
	Name       string `gorm:"column:name;type:varchar(200)" json:"name"`
	Code       string `gorm:"column:code;type:varchar(20)" json:"code"`
	RegionCode string `gorm:"column:region_code;type:varchar(20)" json:"region_code"`
	MainBaseModel
}

type SubRegionHistory struct {
	SubRegion
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
