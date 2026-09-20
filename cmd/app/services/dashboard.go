package services

import (
	"admin-dashboard/cmd/app/repositories"
)

type DashboardData struct {
	Stats            repositories.DashboardStats           `json:"stats"`
	RecentActivities []repositories.RecentActivityItem     `json:"recent_activities"`
	RecentLogins     []repositories.RecentLoginItem        `json:"recent_logins"`
}

func (s *Service) GetDashboardData() (DashboardData, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	stats, err := repository.GetDashboardStats()
	if err != nil {
		if s.Logger != nil {
			s.Logger.Warnf("Failed to retrieve dashboard stats: %v", err)
		}
	}

	activities, err := repository.GetRecentActivities(6)
	if err != nil {
		if s.Logger != nil {
			s.Logger.Warnf("Failed to retrieve recent activities: %v", err)
		}
	}

	logins, err := repository.GetRecentLogins(6)
	if err != nil {
		if s.Logger != nil {
			s.Logger.Warnf("Failed to retrieve recent logins: %v", err)
		}
	}

	return DashboardData{
		Stats:            stats,
		RecentActivities: activities,
		RecentLogins:     logins,
	}, nil
}
