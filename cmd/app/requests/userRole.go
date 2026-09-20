package requests

type UserRoleRequest struct {
	PasswordMinLength  uint64 `param:"password_min_length" json:"password_min_length" form:"password_min_length" query:"password_min_length" validate:"required,min=8,max=20"`
	PasswordMaxLength  uint64 `param:"password_max_length" json:"password_max_length" form:"password_max_length" query:"password_max_length" validate:"required,min=8,max=20,gtefield=PasswordMinLength"`
	SessionTime        uint64  `param:"session_time" json:"session_time" form:"session_time" query:"session_time" validate:"required,min=1,max=100"`
	SessionCheckTime   uint64  `param:"session_check_time" json:"session_check_time" form:"session_check_time" query:"session_check_time" validate:"required,min=1,max=30"`
	PasswordExpireDays uint64  `param:"password_expire_days" json:"password_expire_days" form:"password_expire_days" query:"password_expire_days" validate:"required,min=30,max=365"`
}
