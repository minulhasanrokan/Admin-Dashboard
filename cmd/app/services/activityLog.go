package services

import (
	"admin-dashboard/cmd/app/repositories"
)

type FormattedActivityItem struct {
	Id             uint64 `json:"id"`
	Activity       string `json:"activity"`
	ActivityStatus string `json:"activity_status"`
	TableName      string `json:"table_name"`
	TableId        uint64 `json:"table_id"`
	IPAddress      string `json:"ip_address"`
	Device         string `json:"device"`
	CreatedAt      string `json:"created_at"`
}

type FormattedLoginHistoryItem struct {
	Id           uint64 `json:"id"`
	UserId       int    `json:"user_id"`
	IpAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	LogInTime    string `json:"log_in_time"`
	LogOutStatus int    `json:"log_out_status"`
	StatusText   string `json:"status_text"`
	CreatedAt    string `json:"created_at"`
}

func (s *Service) GetActivityLogList(params repositories.BackupFilterParams) (map[string]interface{}, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	logs, total, err := repository.GetActivityLogs(params)
	if err != nil {
		return nil, err
	}

	rows := make([]FormattedActivityItem, len(logs))
	for i, l := range logs {
		act := ""
		if l.Activity != nil {
			act = *l.Activity
		}
		actStat := ""
		if l.ActivityStatus != nil {
			actStat = *l.ActivityStatus
		}
		tbl := ""
		if l.TableName != nil {
			tbl = *l.TableName
		}
		var tblId uint64
		if l.TableId != nil {
			tblId = *l.TableId
		}
		ip := ""
		if l.IPAddress != nil {
			ip = *l.IPAddress
		}
		dev := ""
		if l.Device != nil {
			dev = *l.Device
		}

		rows[i] = FormattedActivityItem{
			Id:             l.Id,
			Activity:       act,
			ActivityStatus: actStat,
			TableName:      tbl,
			TableId:        tblId,
			IPAddress:      ip,
			Device:         dev,
			CreatedAt:      l.CreatedAt.Format("02 Jan, 2006 15:04:05"),
		}
	}

	result := map[string]interface{}{
		"data":  rows,
		"total": total,
		"page":  params.Page,
		"limit": params.Limit,
	}

	return result, nil
}

func (s *Service) GetLoginHistoryList(params repositories.BackupFilterParams) (map[string]interface{}, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	histories, total, err := repository.GetLoginHistories(params)
	if err != nil {
		return nil, err
	}

	rows := make([]FormattedLoginHistoryItem, len(histories))
	for i, h := range histories {
		uid := 0
		if h.UserId != nil {
			uid = *h.UserId
		}
		ip := ""
		if h.IpAddress != nil {
			ip = *h.IpAddress
		}
		ua := ""
		if h.UserAgent != nil {
			ua = *h.UserAgent
		}
		logInTimeStr := ""
		if h.LogInTime != nil {
			logInTimeStr = h.LogInTime.Format("02 Jan, 2006 15:04:05")
		}
		logOutStat := 0
		if h.LogOutStatus != nil {
			logOutStat = *h.LogOutStatus
		}
		statusText := "Active"
		if logOutStat != 0 {
			statusText = "Logged Out"
		}

		rows[i] = FormattedLoginHistoryItem{
			Id:           h.Id,
			UserId:       uid,
			IpAddress:    ip,
			UserAgent:    ua,
			LogInTime:    logInTimeStr,
			LogOutStatus: logOutStat,
			StatusText:   statusText,
			CreatedAt:    h.CreatedAt.Format("02 Jan, 2006 15:04:05"),
		}
	}

	result := map[string]interface{}{
		"data":  rows,
		"total": total,
		"page":  params.Page,
		"limit": params.Limit,
	}

	return result, nil
}
