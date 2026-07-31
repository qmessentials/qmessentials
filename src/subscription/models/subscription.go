package models

import "time"

type Subscription struct {
	ID          int       `json:"id"`
	OwnerUserID string    `json:"ownerUserId"`
	RuleText    string    `json:"ruleText"`
	VersionID   int       `json:"versionId"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateSubscriptionInput struct {
	RuleText string `json:"ruleText"`
}

type UpdateSubscriptionInput struct {
	RuleText string `json:"ruleText"`
}
