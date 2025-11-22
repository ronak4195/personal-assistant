package models

import "time"

type Category struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"userId" json:"userId"`
	Name      string    `bson:"name" json:"name"`
	ParentID  *string   `bson:"parentId,omitempty" json:"parentId,omitempty"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
