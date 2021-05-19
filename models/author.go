package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthorModel struct {
	ID			primitive.ObjectID		`bson:"_id"`
	LastName	string					`bson:"last_name"`
	FirstName	string					`bson:"first_name"`
	MiddleName	string					`bson:"middle_name"`
}