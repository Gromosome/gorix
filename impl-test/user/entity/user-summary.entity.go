package entity

type UserSummary struct {
	TotalUsers  int64 `json:"totalUsers" db:"total_users"`
	ActiveUsers int64 `json:"activeUsers" db:"active_users"`
}
