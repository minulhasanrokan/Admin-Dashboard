package models

type SystemSetting struct {
	Name             *string `gorm:"column:name;type:varchar(200)" json:"name"`
	Title            *string `gorm:"column:title;type:varchar(200)" json:"title"`
	Slogan           *string `gorm:"column:slogan;type:varchar(200)" json:"slogan"`
	Email            *string `gorm:"column:email;type:varchar(100)" json:"email"`
	Fax              *string `gorm:"column:fax;type:varchar(200)" json:"fax"`
	Mobile           *string `gorm:"column:mobile;type:varchar(20)" json:"mobile"`
	Address          *string `gorm:"column:address;type:text" json:"address"`
	Details          *string `gorm:"column:details;type:text" json:"details"`
	Country          *string `gorm:"column:country;type:varchar(20)" json:"country"`
	State            *string `gorm:"column:state;type:varchar(20)" json:"state"`
	City             *string `gorm:"column:city;type:varchar(20)" json:"city"`
	PostCode         *string `gorm:"column:post_code;type:varchar(20)" json:"post_code"`
	Icon             *string `gorm:"column:icon;type:varchar(200)" json:"icon"`
	AppleIcon        *string `gorm:"column:apple_icon;type:varchar(200)" json:"apple_icon"`
	Version          *string `gorm:"column:version;type:varchar(200)" json:"version"`
	CopyRight        *string `gorm:"column:copy_right;type:varchar(200)" json:"copy_right"`
	LoginBg          *string `gorm:"column:login_bg;type:varchar(200)" json:"login_bg"`
	RegistrationBg   *string `gorm:"column:registration_bg;type:varchar(200)" json:"registration_bg"`
	ForgotPasswordBg *string `gorm:"column:forgot_password_bg;type:varchar(200)" json:"forgot_password_bg"`
	Logo             *string `gorm:"column:logo;type:varchar(200)" json:"logo"`
	LogoWhite        *string `gorm:"column:logo_white;type:varchar(200)" json:"logo_white"`
	LogoSmall        *string `gorm:"column:logo_small;type:varchar(200)" json:"logo_small"`
	LogoDarkSmall    *string `gorm:"column:logo_dark_small;type:varchar(200)" json:"logo_dark_small"`
	VerificationBg   *string `gorm:"column:verification_bg;type:varchar(200)" json:"verification_bg"`
	MainBaseModel
}
type SystemSettingHistory struct {
	SystemSetting
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
