package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ReportID          string             `bson:"report_id" json:"report_id"`
	UserID            *int               `bson:"user_id,omitempty" json:"user_id,omitempty"`
	ClientGeneratedID string             `bson:"client_generated_id" json:"client_generated_id"`
	IsPurchased       bool               `bson:"is_purchased" json:"is_purchased"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
}
