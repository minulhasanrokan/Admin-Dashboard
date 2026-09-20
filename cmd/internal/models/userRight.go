package models

type UserRight struct {
	UserId   uint64   `gorm:"column:user_id;type:int unsigned" json:"user_id"`
	RightKey string `gorm:"column:right_key;type:varchar(255)" json:"right_key"`
	MainBaseModel
}

type UserRightHistory struct {
	UserRight
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
