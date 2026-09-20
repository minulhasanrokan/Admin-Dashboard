package requests

import "mime/multipart"

type BusinessSettingRequest struct {
	Name      string `param:"name" json:"name" form:"name" query:"name" validate:"required,min=5,max=200"`
	Title     string `param:"title" json:"title" form:"title" query:"title" validate:"required,min=5,max=200"`
	Slogan    string `param:"slogan" json:"slogan" form:"slogan" query:"slogan" validate:"required,min=5,max=200"`
	Email     string `param:"email" json:"email" form:"email" query:"email" validate:"required,email"`
	Mobile    string `param:"mobile" json:"mobile" form:"mobile" query:"mobile" validate:"required,mobile"`
	Fax       string `param:"fax" json:"fax" form:"fax" query:"fax" validate:"required,min=5,max=200"`
	Address   string `param:"address" json:"address" form:"address" query:"address" validate:"required,min=5,max=200"`
	Country   string `param:"country" json:"country" form:"country" query:"country" validate:"required,max=20"`
	State     string `param:"state" json:"state" form:"state" query:"state" validate:"required,max=20"`
	City      string `param:"city" json:"city" form:"city" query:"city" validate:"required,max=20"`
	PostCode  string `param:"post_code" json:"post_code" form:"post_code" query:"post_code" validate:"required,min=4,max=4"`
	Details   string `param:"details" json:"details" form:"details" query:"details" validate:"omitempty,max=200"`
	Version   string `param:"version" json:"version" form:"version" query:"version" validate:"required,min=1,max=200"`
	CopyRight string `param:"copy_right" json:"copy_right" form:"copy_right" query:"copy_right" validate:"required,min=5,max=200"`
}

type ImagesSettingRequest struct {
	Logo             *multipart.FileHeader `param:"logo" json:"logo" form:"logo" query:"logo" validate:"omitempty,file,minsize=1,maxsize_mb=5,imagetype"`
	LogoWhite        *multipart.FileHeader `param:"logo_white" json:"logo_white" form:"logo_white" query:"logo_white" validate:"omitempty,file,minsize=1,maxsize_mb=5,imagetype"`
	LogoSmall        *multipart.FileHeader `param:"logo_small" json:"logo_small" form:"logo_small" query:"logo_small" validate:"omitempty,file,minsize=1,maxsize_mb=5,imagetype"`
	LogoDarkSmall    *multipart.FileHeader `param:"logo_dark_small" json:"logo_dark_small" form:"logo_dark_small" query:"logo_dark_small" validate:"omitempty,file,minsize=1,maxsize_mb=5,imagetype"`
	Icon             *multipart.FileHeader `param:"icon" json:"icon" form:"icon" query:"icon" validate:"omitempty,file,minsize=1,maxsize_mb=5,imagetype"`
	AppleIcon        *multipart.FileHeader `param:"apple_icon" json:"apple_icon" form:"apple_icon" query:"apple_icon" validate:"omitempty,file,minsize=1,maxsize_mb=5,imagetype"`
	LoginBg          *multipart.FileHeader `param:"login_bg" json:"login_bg" form:"login_bg" query:"login_bg" validate:"omitempty,file,minsize=1,maxsize_mb=5,imagetype"`
	ForgotPasswordBg *multipart.FileHeader `param:"forgot_password_bg" json:"forgot_password_bg" form:"forgot_password_bg" query:"forgot_password_bg" validate:"omitempty,file,minsize=1,maxsize_mb=5,imagetype"`
	RegistrationBg   *multipart.FileHeader `param:"registration_bg" json:"registration_bg" form:"registration_bg" query:"registration_bg" validate:"omitempty,file,minsize=1,maxsize_mb=5,imagetype"`
	VerificationBg   *multipart.FileHeader `param:"verification_bg" json:"verification_bg" form:"verification_bg" query:"verification_bg" validate:"omitempty,file,minsize=1,maxsize_mb=5,imagetype"`
}

type BackupSettingRequest struct {
	Backup string `param:"backup" json:"backup" form:"backup" query:"backup" validate:"required,oneof=0001 0002"`
}
