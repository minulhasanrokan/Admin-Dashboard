package models

type Designation struct {
	Name    *string `gorm:"column:name;type:varchar(200)" json:"name"`
	Title   *string `gorm:"column:title;type:varchar(250)" json:"title"`
	Code    *string `gorm:"column:code;type:varchar(50)" json:"code"`
	Details *string `gorm:"column:details;type:text" json:"details"`
	MainBaseModel
}

type DesignationHistory struct {
	Designation
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
