package models

type Country struct {
	Name          *string `gorm:"column:name;type:varchar(200)" json:"name"`
	Alpha2        *string `gorm:"column:alpha_2;type:varchar(2)" json:"alpha_2"`
	Alpha3        *string `gorm:"column:alpha_3;type:varchar(3)" json:"alpha_3"`
	Code          *string `gorm:"column:code;type:varchar(20)" json:"code"`
	Iso31662      *string `gorm:"column:iso_3166_2;type:varchar(20)" json:"iso_3166_2"`
	RegionCode    *string `gorm:"column:region_code;type:varchar(20)" json:"region_code"`
	SubRegionCode *string `gorm:"column:sub_region_code;type:varchar(20)" json:"sub_region_code"`
	MainBaseModel
}

type CountryHistory struct {
	Country
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
