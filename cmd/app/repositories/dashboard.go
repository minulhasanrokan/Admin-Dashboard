package repositories

import (
	"admin-dashboard/cmd/internal/models"
	"time"
)

type DashboardStats struct {
	TotalUsers         int64 `json:"total_users"`
	ActiveUsers        int64 `json:"active_users"`
	InactiveUsers      int64 `json:"inactive_users"`
	VerifiedUsers      int64 `json:"verified_users"`
	NewUsersThisMonth  int64 `json:"new_users_this_month"`
	NewUsersThisWeek   int64 `json:"new_users_this_week"`
	TotalBackups       int64 `json:"total_backups"`
	SuccessfulBackups  int64 `json:"successful_backups"`
	TotalEmailsSent    int64 `json:"total_emails_sent"`
	PendingEmails      int64 `json:"pending_emails"`
	TotalSmsSent       int64 `json:"total_sms_sent"`
	PendingSms         int64 `json:"pending_sms"`
}

type RecentActivityItem struct {
	Id         uint64    `json:"id"`
	Activity   string    `json:"activity"`
	TableName  string    `json:"table_name"`
	TableId    uint64    `json:"table_id"`
	IPAddress  string    `json:"ip_address"`
	Device     string    `json:"device"`
	CreatedAt  time.Time `json:"created_at"`
}

type RecentLoginItem struct {
	Id             uint64     `json:"id"`
	UserId         int        `json:"user_id"`
	UserName       string     `json:"user_name"`
	IpAddress      string     `json:"ip_address"`
	UserAgent      string     `json:"user_agent"`
	LogInTime      *time.Time `json:"log_in_time"`
	LogOutStatus   int        `json:"log_out_status"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (r *Repository) GetDashboardStats() (DashboardStats, error) {
	var stats DashboardStats

	if r.Db == nil || r.Db.GetDb() == nil {
		return stats, nil
	}

	db := r.Db.GetDb()

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())

	// User Stats
	_ = db.Model(&models.User{}).Where("delete_status = ?", 0).Count(&stats.TotalUsers)
	_ = db.Model(&models.User{}).Where("delete_status = ? AND status = ?", 0, 1).Count(&stats.ActiveUsers)
	_ = db.Model(&models.User{}).Where("delete_status = ? AND status = ?", 0, 0).Count(&stats.InactiveUsers)
	_ = db.Model(&models.User{}).Where("delete_status = ? AND verified = ?", 0, 1).Count(&stats.VerifiedUsers)
	_ = db.Model(&models.User{}).Where("delete_status = ? AND created_at >= ?", 0, startOfMonth).Count(&stats.NewUsersThisMonth)
	_ = db.Model(&models.User{}).Where("delete_status = ? AND created_at >= ?", 0, startOfWeek).Count(&stats.NewUsersThisWeek)

	// Backup Stats
	_ = db.Model(&models.BackupHistory{}).Count(&stats.TotalBackups)
	_ = db.Model(&models.BackupHistory{}).Where("backup_status = ?", 1).Count(&stats.SuccessfulBackups)

	// Mail Stats
	_ = db.Model(&models.MailHistory{}).Where("sent_status = ?", 1).Count(&stats.TotalEmailsSent)
	_ = db.Model(&models.MailHistory{}).Where("sent_status = ?", 0).Count(&stats.PendingEmails)

	// SMS Stats
	_ = db.Model(&models.SmsHistory{}).Where("sent_status = ?", 1).Count(&stats.TotalSmsSent)
	_ = db.Model(&models.SmsHistory{}).Where("sent_status = ?", 0).Count(&stats.PendingSms)

	return stats, nil
}

func (r *Repository) GetRecentActivities(limit int) ([]RecentActivityItem, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return []RecentActivityItem{}, nil
	}

	if limit <= 0 {
		limit = 6
	}

	var rawActivities []models.UserActivity
	err := r.Db.GetDb().Model(&models.UserActivity{}).
		Order("created_at DESC").
		Limit(limit).
		Find(&rawActivities).Error

	if err != nil {
		return nil, err
	}

	items := make([]RecentActivityItem, len(rawActivities))
	for i, a := range rawActivities {
		act := ""
		if a.Activity != nil {
			act = *a.Activity
		}
		tbl := ""
		if a.TableName != nil {
			tbl = *a.TableName
		}
		var tblId uint64
		if a.TableId != nil {
			tblId = *a.TableId
		}
		ip := ""
		if a.IPAddress != nil {
			ip = *a.IPAddress
		}
		dev := ""
		if a.Device != nil {
			dev = *a.Device
		}

		items[i] = RecentActivityItem{
			Id:        a.Id,
			Activity:  act,
			TableName: tbl,
			TableId:   tblId,
			IPAddress: ip,
			Device:    dev,
			CreatedAt: a.CreatedAt,
		}
	}

	return items, nil
}

func (r *Repository) GetRecentLogins(limit int) ([]RecentLoginItem, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return []RecentLoginItem{}, nil
	}

	if limit <= 0 {
		limit = 6
	}

	var rawLogins []models.UserLoginHistory
	err := r.Db.GetDb().Model(&models.UserLoginHistory{}).
		Order("created_at DESC").
		Limit(limit).
		Find(&rawLogins).Error

	if err != nil {
		return nil, err
	}

	items := make([]RecentLoginItem, len(rawLogins))
	for i, l := range rawLogins {
		uid := 0
		if l.UserId != nil {
			uid = *l.UserId
		}
		ip := ""
		if l.IpAddress != nil {
			ip = *l.IpAddress
		}
		ua := ""
		if l.UserAgent != nil {
			ua = *l.UserAgent
		}
		logOut := 0
		if l.LogOutStatus != nil {
			logOut = *l.LogOutStatus
		}

		// Optionally find username
		userName := ""
		if uid > 0 {
			var u models.User
			if err := r.Db.GetDb().Select("user_name, name").Where("id = ?", uid).First(&u).Error; err == nil {
				if u.Name != nil && *u.Name != "" {
					userName = *u.Name
				} else if u.UserName != nil {
					userName = *u.UserName
				}
			}
		}

		items[i] = RecentLoginItem{
			Id:           l.Id,
			UserId:       uid,
			UserName:     userName,
			IpAddress:    ip,
			UserAgent:    ua,
			LogInTime:    l.LogInTime,
			LogOutStatus: logOut,
			CreatedAt:    l.CreatedAt,
		}
	}

	return items, nil
}
