package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
) 

type Donation struct {
    ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
    Amount    float64            `json:"amount" bson:"amount"`
    Currency  string             `json:"currency" bson:"currency"`
    Name      string             `json:"name" bson:"name"`
    Message   string             `json:"message" bson:"message"`
    Streamer  string             `json:"streamer" bson:"streamer"`
    Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
}