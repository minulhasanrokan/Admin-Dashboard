package models

type PrivacyPolicy struct {
	Content *string `gorm:"column:content;type:text" json:"content"`
	MainBaseModel
}

type PrivacyPolicyHistory struct {
	PrivacyPolicy
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
