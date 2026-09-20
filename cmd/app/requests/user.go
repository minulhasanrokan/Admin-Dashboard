package requests

type UserRegisterRequest struct {
	Name            string `param:"name" json:"name" form:"name" query:"name" validate:"required,min=3,max=100"`
	Email           string `param:"email" json:"email" form:"email" query:"email" validate:"required,email"`
	Mobile          string `param:"mobile" json:"mobile" form:"mobile" query:"mobile" validate:"required,mobile"`
	Password        string `param:"password" json:"password" form:"password" query:"password" validate:"omitempty,password"`
	ConfirmPassword string `param:"confirm_password" json:"confirm_password" form:"confirm_password" query:"confirm_password" validate:"omitempty,eqfield=Password"`
	TermsPrivacy    int    `param:"terms_privacy" json:"terms_privacy" form:"terms_privacy" query:"terms_privacy" validate:"required,eq=1"`
	InviteCode      string `param:"invite_code" json:"invite_code" form:"invite_code" query:"invite_code"`
}

type UserUUIDRequest struct {
	UUID string `param:"uuid" json:"uuid" form:"uuid" query:"uuid" validate:"required,uuid4"`
}

type UserVerifyEmailRequest struct {
	UUID       string `param:"uuid" json:"uuid" form:"uuid" query:"uuid" validate:"required,uuid4"`
	DigitOne   string `param:"digit_one" json:"digit_one" form:"digit_one" query:"digit_one" validate:"required,numeric,len=1"`
	DigitTwo   string `param:"digit_two" json:"digit_two" form:"digit_two" query:"digit_two" validate:"required,numeric,len=1"`
	DigitThree string `param:"digit_three" json:"digit_three" form:"digit_three" query:"digit_three" validate:"required,numeric,len=1"`
	DigitFour  string `param:"digit_four" json:"digit_four" form:"digit_four" query:"digit_four" validate:"required,numeric,len=1"`
}

type UserPasswordResetRequest struct {
	UUID            string `param:"uuid" json:"uuid" form:"uuid" query:"uuid" validate:"required,uuid4"`
	DigitOne        string `param:"digit_one" json:"digit_one" form:"digit_one" query:"digit_one" validate:"required,numeric,len=1"`
	DigitTwo        string `param:"digit_two" json:"digit_two" form:"digit_two" query:"digit_two" validate:"required,numeric,len=1"`
	DigitThree      string `param:"digit_three" json:"digit_three" form:"digit_three" query:"digit_three" validate:"required,numeric,len=1"`
	DigitFour       string `param:"digit_four" json:"digit_four" form:"digit_four" query:"digit_four" validate:"required,numeric,len=1"`
	Password        string `param:"password" json:"password" form:"password" query:"password" validate:"required,password"`
	ConfirmPassword string `param:"confirm_password" json:"confirm_password" form:"confirm_password" query:"confirm_password" validate:"required,eqfield=Password"`
}

type EmailRequest struct {
	Email string `param:"email" json:"email" form:"email" query:"email" validate:"required,email"`
}

type PasswordResetQuery struct {
	RawEmail string `query:"email" param:"email" json:"email" form:"email" validate:"required"`
	Token    string `query:"token" param:"token" json:"token" form:"token" validate:"required"`
}
