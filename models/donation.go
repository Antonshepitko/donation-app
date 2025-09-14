package models

import "time"

type Donation struct {
	Streamer  string    `json:"streamer" bson:"streamer"`
	DonorName string    `json:"donorName" bson:"donorName"`
	Amount    float64   `json:"amount" bson:"amount"`
	Currency  string    `json:"currency" bson:"currency"`
	Message   string    `json:"message" bson:"message"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type Donation struct {
    ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
    Amount    float64            `json:"amount" bson:"amount"`
    Currency  string             `json:"currency" bson:"currency"`
    Name      string             `json:"name" bson:"name"`
    Message   string             `json:"message" bson:"message"`
    Streamer  string             `json:"streamer" bson:"streamer"`
    Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
}