package constants

const (
	ActivityAdd           = 1
	ActivityAddEdit       = 2
	ActivityEdit          = 3
	ActivityDelete        = 4
	ActivityDeactivate    = 5
	ActivityActivate      = 6
	ActivityResetPassword = 7
	ActivitySetMenuRight  = 8
	ActivitySetFileRight  = 9
	ActivityUnlock        = 10
	ActivitySentToChecker = 11
	ActivityApprove       = 12
	ActivityReturn        = 13
	ActivityDownload      = 14
	ActivityReport        = 15
	ActivityReceive       = 16
	ActivityStatusUpdate  = 17
	ActivityUpload        = 18
	ActivitySendEmail     = 19
	ActivitySendSMS       = 20
)

var ActivityList = map[int64]string{
	ActivityAdd:           "New Add",
	ActivityAddEdit:       "Add-Edit",
	ActivityEdit:          "Edit",
	ActivityDelete:        "Delete",
	ActivityDeactivate:    "Deactivate",
	ActivityActivate:      "Active",
	ActivityResetPassword: "Reset Password",
	ActivitySetMenuRight:  "Set Menu Right",
	ActivitySetFileRight:  "Set File Right",
	ActivityUnlock:        "Unlock",
	ActivitySentToChecker: "Sent to Checker",
	ActivityApprove:       "Approve",
	ActivityReturn:        "Return",
	ActivityDownload:      "Download",
	ActivityReport:        "Report",
	ActivityReceive:       "Receive",
	ActivityStatusUpdate:  "Status Update",
	ActivityUpload:        "Upload",
	ActivitySendEmail:     "Send Email",
	ActivitySendSMS:       "Send SMS",
}

const (
	ProcessStatusNew           = 0
	ProcessStatusSentToChecker = 1
	ProcessStatusReturned      = 2
	ProcessStatusApproved      = 3
)

var ProcessStatusList = map[int64]string{
	ProcessStatusNew:           "",
	ProcessStatusSentToChecker: "(Sent to Checker)",
	ProcessStatusReturned:      "(Returned)",
	ProcessStatusApproved:      "(Approved)",
}

const (
	DataStatusDeactivate = 0
	DataStatusActive     = 1
)

var DataActiveStatusList = map[int64]string{
	DataStatusDeactivate: "Deactivate",
	DataStatusActive:     "Active",
}
