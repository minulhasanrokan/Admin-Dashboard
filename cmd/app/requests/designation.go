package requests

type CreateDesignationRequest struct {
	Name            string `param:"name" json:"name" form:"name" query:"name"`
	DesignationName string `param:"designation_name" json:"designation_name" form:"designation_name" query:"designation_name"`
	Title           string `param:"title" json:"title" form:"title" query:"title"`
	DesignationTitle string `param:"designation_title" json:"designation_title" form:"designation_title" query:"designation_title"`
	Code            string `param:"code" json:"code" form:"code" query:"code"`
	DesignationCode string `param:"designation_code" json:"designation_code" form:"designation_code" query:"designation_code"`
	Details         string `param:"details" json:"details" form:"details" query:"details"`
	DesignationDetails string `param:"designation_details" json:"designation_details" form:"designation_details" query:"designation_details"`
}

func (r *CreateDesignationRequest) GetName() string {
	if r.DesignationName != "" {
		return r.DesignationName
	}
	return r.Name
}

func (r *CreateDesignationRequest) GetTitle() string {
	if r.DesignationTitle != "" {
		return r.DesignationTitle
	}
	return r.Title
}

func (r *CreateDesignationRequest) GetCode() string {
	if r.DesignationCode != "" {
		return r.DesignationCode
	}
	return r.Code
}

func (r *CreateDesignationRequest) GetDetails() string {
	if r.DesignationDetails != "" {
		return r.DesignationDetails
	}
	return r.Details
}

type UpdateDesignationRequest struct {
	UpdateId        uint64 `param:"update_id" json:"update_id" form:"update_id" query:"update_id"`
	Id              uint64 `param:"id" json:"id" form:"id" query:"id"`
	Name            string `param:"name" json:"name" form:"name" query:"name"`
	DesignationName string `param:"designation_name" json:"designation_name" form:"designation_name" query:"designation_name"`
	Title           string `param:"title" json:"title" form:"title" query:"title"`
	DesignationTitle string `param:"designation_title" json:"designation_title" form:"designation_title" query:"designation_title"`
	Code            string `param:"code" json:"code" form:"code" query:"code"`
	DesignationCode string `param:"designation_code" json:"designation_code" form:"designation_code" query:"designation_code"`
	Details         string `param:"details" json:"details" form:"details" query:"details"`
	DesignationDetails string `param:"designation_details" json:"designation_details" form:"designation_details" query:"designation_details"`
}

func (r *UpdateDesignationRequest) GetName() string {
	if r.DesignationName != "" {
		return r.DesignationName
	}
	return r.Name
}

func (r *UpdateDesignationRequest) GetTitle() string {
	if r.DesignationTitle != "" {
		return r.DesignationTitle
	}
	return r.Title
}

func (r *UpdateDesignationRequest) GetCode() string {
	if r.DesignationCode != "" {
		return r.DesignationCode
	}
	return r.Code
}

func (r *UpdateDesignationRequest) GetDetails() string {
	if r.DesignationDetails != "" {
		return r.DesignationDetails
	}
	return r.Details
}

type DeleteDesignationRequest struct {
	DeleteReason string `param:"delete_reason" json:"delete_reason" form:"delete_reason" query:"delete_reason" validate:"required,max=250"`
}

type SentToCheckerDesignationRequest struct {
	Checker    uint64 `param:"checker" json:"checker" form:"checker" query:"checker" validate:"required"`
	NotifyType uint64 `param:"notify_type" json:"notify_type" form:"notify_type" query:"notify_type" validate:"required"`
	Remarks    string `param:"remarks" json:"remarks" form:"remarks" query:"remarks" validate:"max=250"`
}

type VerifyDesignationRequest struct {
	VerifyType   int    `param:"verify_type" json:"verify_type" form:"verify_type" query:"verify_type"`
	ReturnReason string `param:"return_reason" json:"return_reason" form:"return_reason" query:"return_reason"`
}

type StatusReasonDesignationRequest struct {
	Reason         string `param:"reason" json:"reason" form:"reason" query:"reason"`
	ActiveReason   string `param:"active_reason" json:"active_reason" form:"active_reason" query:"active_reason"`
	DeactiveReason string `param:"deactive_reason" json:"deactive_reason" form:"deactive_reason" query:"deactive_reason"`
}

func (r *StatusReasonDesignationRequest) GetReason() string {
	if r.ActiveReason != "" {
		return r.ActiveReason
	}
	if r.DeactiveReason != "" {
		return r.DeactiveReason
	}
	return r.Reason
}
