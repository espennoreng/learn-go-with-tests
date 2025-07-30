package models

import "time"

type Item struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	ExternalID string    `json:"externalId"`
	OrgID      string    `json:"orgId"`
	IsActive   string    `json:"isActive"`
	CreatedAt  time.Time `json:"createdAt"`
	CreatedBy  string    `json:"createdBy"`
	DeletedAt  string    `json:"deletedAt"`
}
