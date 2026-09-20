package seed

import (
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/models"
)

func BaseSeedData() models.MainBaseModel {
	return models.MainBaseModel{
		Status:        1,
		DeleteStatus:  0,
		ProcessStatus: *helpers.UInt(0),
		ActualStatus:  *helpers.UInt(0),
		EditStatus:    0,
		ActiveStatus:  *helpers.UInt(1),
		Checker:       *helpers.UInt(0),
		NotifyType:    *helpers.UInt(0),
		Message:       "",
		Remarks:       "",
	}
}
