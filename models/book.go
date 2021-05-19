package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BookModel struct {
	ID			primitive.ObjectID		`bson:"_id"`
	BookId		string					`bson:"id"`
	BookName	string					`bson:"name"`
	Authors		[]primitive.ObjectID	`bson:"authors"`
	Series		string					`bson:"series"`
	SerialNum	string					`bson:"sernum"`
	Filename	string					`bson:"filename"`
	Deleted		bool					`bson:"deleted"`
	Language	string					`bson:"lang"`
	Keywords	[]string				`bson:"keywords"`
	Added		string					`bson:"added"`
	Genres		[]string				`bson:"genres"`
	Year		string					`bson:"year"`
	ISBN		string					`bson:"isbn"`
	City		string					`bson:"city"`
	PubName		string					`bson:"pub_name"`
	Publisher	string					`bson:"publisher"`
	Rate		string					`bson:"rate"`
}
