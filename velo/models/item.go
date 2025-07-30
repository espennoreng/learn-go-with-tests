package models

import "time"

type Item struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	ExternalID string    `json:"external_id"`
	OrgID      string    `json:"org_id"`
	IsActive   string    `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	DeletedAt  string    `json:"deleted_at"`
}
