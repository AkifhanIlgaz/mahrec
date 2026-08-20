package models

import "time"

type Goal struct {
	ID           int64
	ClientID     string
	Title        string
	Description  string
	TargetCount  int
	CurrentCount int
	CreatedAt    time.Time
	CompletedAt  *time.Time
}

func (g Goal) IsCompleted() bool {
	return g.CompletedAt != nil
}

func (g Goal) Progress() int {
	if g.TargetCount <= 0 {
		return 0
	}
	pct := g.CurrentCount * 100 / g.TargetCount
	if pct > 100 {
		return 100
	}
	return pct
}
