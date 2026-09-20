package models

import "time"

type User struct {
	Name                   *string    `gorm:"column:name;type:varchar(200)" json:"name"`
	UserName               *string    `gorm:"column:user_name;type:varchar(250)" json:"user_name"`
	UserId                 *string    `gorm:"column:user_id;type:varchar(50)" json:"user_id"`
	Email                  *string    `gorm:"column:email;type:varchar(100)" json:"email"`
	Mobile                 *string    `gorm:"column:mobile;type:varchar(20)" json:"mobile"`
	Address                *string    `gorm:"column:address;type:varchar(200)" json:"address"`
	Details                *string    `gorm:"column:details;type:text" json:"details"`
	Department             *uint64      `gorm:"column:department;type:int unsigned" json:"department"`
	AssignedDepartment     *string    `gorm:"column:assigned_department;type:varchar(100)" json:"assigned_department"`
	Designation            *uint64      `gorm:"column:designation;type:int unsigned" json:"designation"`
	LineManagerStatus      *uint64      `gorm:"column:line_manager_status;type:int unsigned" json:"line_manager_status"`
	FirstLineManager       *uint64      `gorm:"column:first_line_manager;type:int unsigned" json:"first_line_manager"`
	SecondLineManager      *uint64      `gorm:"column:second_line_manager;type:int unsigned" json:"second_line_manager"`
	Group                  *uint64      `gorm:"column:group;type:int unsigned" json:"group"`
	BirthDate              *time.Time `gorm:"column:birth_date;type:date" json:"birth_date"`
	Gender                 *uint64      `gorm:"column:gender;type:int unsigned" json:"gender"`
	BloodGroup             *uint64      `gorm:"column:blood_group;type:int unsigned" json:"blood_group"`
	Photo                  *string    `gorm:"column:photo;type:varchar(200)" json:"photo"`
	SessionTime            *uint64      `gorm:"column:session_time;type:int unsigned" json:"session_time"`
	SessionCheckTime       *uint64      `gorm:"column:session_check_time;type:int unsigned" json:"session_check_time"`
	SuperAdminStatus       uint64       `gorm:"column:super_admin_status;default:0" json:"super_admin_status"`
	LockStatus             *uint64      `gorm:"column:lock_status;type:int unsigned" json:"lock_status"`
	BlockStatus            *uint64      `gorm:"column:block_status;type:int unsigned" json:"block_status"`
	UserType               uint64       `gorm:"column:user_type;type:int unsigned" json:"user_type"`
	Password               *string    `gorm:"column:password;type:varchar(250)" json:"-"`
	PasswordChangeStatus   *uint64      `gorm:"column:password_change_status;type:int unsigned" json:"password_change_status"`
	PasswordChangeDate     *time.Time `gorm:"column:password_change_date;type:date" json:"password_change_date"`
	PasswordExpireDate     *time.Time `gorm:"column:password_expire_date;type:date" json:"password_expire_date"`
	LastPasswordChangeDate *time.Time `gorm:"column:last_password_change_date;type:date" json:"last_password_change_date,omitempty"`
	BranchId               *uint64      `gorm:"column:branch_id;type:int unsigned" json:"branch_id"`
	RememberToken          *string    `gorm:"column:remember_token;type:varchar(200)" json:"remember_token"`
	Verified               bool       `gorm:"column:verified;default:false" json:"verified"`
	SignupIP               string     `gorm:"column:signup_ip;type:varchar(50)" json:"signup_ip,omitempty"`
	IsActive               bool       `gorm:"column:is_active;default:true" json:"is_active"`
	TermsPrivacy           int        `gorm:"column:terms_privacy;default:0" json:"terms_privacy"`
	LastLoginIP            string     `gorm:"column:last_login_ip;type:varchar(50)" json:"last_login_ip,omitempty"`
	LastLogin              *time.Time `gorm:"column:last_login" json:"last_login,omitempty"`
	MainBaseModel
}

type UserHistory struct {
	User
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
