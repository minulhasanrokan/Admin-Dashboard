package requests

type AuthenticationSettingRequest struct {
	AllowRegistration    bool   `param:"allow_registration" json:"allow_registration" form:"allow_registration" query:"allow_registration"`
	InviteOnly           bool   `param:"invite_only" json:"invite_only" form:"invite_only" query:"invite_only"`
	VerificationRequired bool   `param:"verification_required" json:"verification_required" form:"verification_required" query:"verification_required"`
	VerificationExpired  uint64 `param:"verification_expired" json:"verification_expired" form:"verification_expired" query:"verification_expired" validate:"required,min=1,max=10080"`
	ReferralSystem       bool   `param:"referral_system" json:"referral_system" form:"referral_system" query:"referral_system"`
	LoginType            string `param:"login_type" json:"login_type" form:"login_type" query:"login_type" validate:"required,oneof=mobile email"`
	PasswordRequired     bool   `param:"password_required" json:"password_required" form:"password_required" query:"password_required"`
	OtpSystem            bool   `param:"otp_system" json:"otp_system" form:"otp_system" query:"otp_system"`
	OtpType              string `param:"otp_type" json:"otp_type" form:"otp_type" query:"otp_type" validate:"required,oneof=sms email"`
}
