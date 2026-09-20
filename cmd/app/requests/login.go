package requests

type UserLoginRequest struct {
	Email    string `param:"email" json:"email" form:"email" query:"email"`
	Mobile   string `param:"mobile" json:"mobile" form:"mobile" query:"mobile"`
	Password string `param:"password" json:"password" form:"password" query:"password"`
}
