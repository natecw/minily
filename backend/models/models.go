package models

import "time"

type CreateRequest struct {
	URL       string     `json:"long_url"`
	Alias     string     `json:"alias,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedBy string     `json:"created_by,omitempty"`
}

type Minily struct {
	ShortCode string `json:"short_url"`
}

type LongUrl struct {
	LongUrl string
}
