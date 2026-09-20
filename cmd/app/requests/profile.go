package requests

import "mime/multipart"

type ProfileUpdateRequest struct {
	Name    string `param:"name" json:"name" form:"name" query:"name" validate:"required,min=3,max=100"`
	Mobile  string `param:"mobile" json:"mobile" form:"mobile" query:"mobile" validate:"omitempty,max=20"`
	Address string `param:"address" json:"address" form:"address" query:"address" validate:"omitempty,max=200"`
	Details string `param:"details" json:"details" form:"details" query:"details" validate:"omitempty,max=1000"`
}

type ProfilePasswordRequest struct {
	CurrentPassword string `param:"current_password" json:"current_password" form:"current_password" query:"current_password" validate:"required,min=6,max=50"`
	NewPassword     string `param:"new_password" json:"new_password" form:"new_password" query:"new_password" validate:"required,password"`
	ConfirmPassword string `param:"confirm_password" json:"confirm_password" form:"confirm_password" query:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type ProfilePhotoRequest struct {
	Photo *multipart.FileHeader `form:"photo" validate:"required,file,minsize=1,maxsize_mb=2,imagetype"`
}
