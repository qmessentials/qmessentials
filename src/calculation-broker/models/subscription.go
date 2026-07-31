package models

type Subscription struct {
	ID        int    `json:"id"`
	RuleText  string `json:"ruleText"`
	VersionID int    `json:"versionId"`
	IsActive  bool   `json:"isActive"`
}
