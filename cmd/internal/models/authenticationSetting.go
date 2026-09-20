package models

type AuthenticationSetting struct {
	AllowRegistration    bool   `gorm:"column:allow_registration;default:true" json:"allow_registration"`
	InviteOnly           bool   `gorm:"column:invite_only;default:false" json:"invite_only"`
	VerificationRequired bool   `gorm:"column:verification_required;default:true" json:"verification_required"`
	VerificationExpired  uint64 `gorm:"column:verification_expired;type:int unsigned;default:15" json:"verification_expired"`
	ReferralSystem       bool   `gorm:"column:referral_system;default:false" json:"referral_system"`
	LoginType            string `gorm:"column:login_type;type:varchar(20);default:'email'" json:"login_type"`
	PasswordRequired     bool   `gorm:"column:password_required;default:true" json:"password_required"`
	OtpSystem            bool   `gorm:"column:otp_system;default:false" json:"otp_system"`
	OtpType              string `gorm:"column:otp_type;type:varchar(20);default:'email'" json:"otp_type"`
	MainBaseModel
}

type AuthenticationSettingHistory struct {
	AuthenticationSetting
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
