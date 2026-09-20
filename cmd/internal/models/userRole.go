package models

type UserRole struct {
	PasswordMinLength  uint64 `gorm:"column:password_min_length;type:tinyint unsigned" json:"password_min_length"`
	PasswordMaxLength  uint64 `gorm:"column:password_max_length;type:tinyint unsigned" json:"password_max_length"`
	SessionTime        uint64  `gorm:"column:session_time;type:int unsigned" json:"session_time"`
	SessionCheckTime   uint64  `gorm:"column:session_check_time;type:int unsigned" json:"session_check_time"`
	PasswordExpireDays uint64  `gorm:"column:password_expire_days;type:int unsigned" json:"password_expire_days"`
	MainBaseModel
}

type UserRoleHistory struct {
	UserRole
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
