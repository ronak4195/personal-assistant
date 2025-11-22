package models

import "time"

type RepeatInterval string

const (
	RepeatNone    RepeatInterval = "none"
	RepeatDaily   RepeatInterval = "daily"
	RepeatWeekly  RepeatInterval = "weekly"
	RepeatMonthly RepeatInterval = "monthly"
)

type Reminder struct {
	ID             string         `bson:"_id,omitempty" json:"id"`
	UserID         string         `bson:"userId" json:"userId"`
	Title          string         `bson:"title" json:"title"`
	Description    *string        `bson:"description,omitempty" json:"description,omitempty"`
	DueAt          time.Time      `bson:"dueAt" json:"dueAt"`
	RepeatInterval RepeatInterval `bson:"repeatInterval" json:"repeatInterval"`
	IsActive       bool           `bson:"isActive" json:"isActive"`
	CreatedAt      time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time      `bson:"updatedAt" json:"updatedAt"`
}
