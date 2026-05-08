package dto

import "time"

type SubmitJobResponse struct {
	JobID string `json:"job_id"`
}

type JobStatusResponse struct {
	JobID      string    `json:"job_id"`
	Status     string    `json:"status"`
	ModelURL   string    `json:"model_url"`
	TextureURL string    `json:"texture_url"`
	TreesURL   string    `json:"trees_url,omitempty"`
	Error      string    `json:"error,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type LimitResponse struct {
	UserID     string    `json:"user_id"`
	DailyLimit int32     `json:"daily_limit"`
	UsedToday  int32     `json:"used_today"`
	ResetAt    time.Time `json:"reset_at"`
}

type SetLimitRequest struct {
	UserID     string `json:"user_id"`
	DailyLimit int32  `json:"daily_limit"`
}
