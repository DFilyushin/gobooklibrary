package database

import (
	"context"
	"github.com/DFilyushin/gobooklibrary/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var DbClient *mongo.Client
var Db *mongo.Database
var DbContext context.Context

const DefaultTimeoutForOperation = 20

func SetupDatabase(conString string, databaseName string)  {
	var err error

	DbContext = context.TODO() //context.WithTimeout(context.Background(), 20*time.Second)
	clientOptions := options.Client().ApplyURI(conString)
	//clientOptions.SetMaxPoolSize(200)
	DbClient, err = mongo.Connect(DbContext, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	//err = DbClient.Ping(DbContext, readpref.Primary())
	//if err != nil {
	//	panic("Mongo unavailable...")
	//}

	Db = DbClient.Database(databaseName)
}

func GetBookById(id string) (*models.BookModel, error) {
	var filter = bson.M{"filename": id}
	var result *models.BookModel
	err := Db.Collection("books").FindOne(DbContext, filter).Decode(&result)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	return result, nil
}

func GetAuthorByFullNames(lastName string, firstName string, middleName string) (*models.AuthorModel, error) {
	var filter = bson.M{"last_name": lastName, "first_name": firstName, "middle_name": middleName}
	var result *models.AuthorModel
	err := Db.Collection("authors").FindOne(DbContext, filter).Decode(&result)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	return result, nil
}

func AddAuthor(author *models.AuthorModel) (*primitive.ObjectID, error)  {
	/*
		Add author to database
	 */
	result, err := Db.Collection("authors").InsertOne(DbContext, &author)
	if err != nil {
		return nil, err
	}
	value := result.InsertedID.(primitive.ObjectID)
	return &value, nil
}

func AddBook(book *models.BookModel) (*primitive.ObjectID, error) {
	/*
		Add book to database
	 */
	result, err := Db.Collection("books").InsertOne(DbContext, &book)
	if err != nil {
		return nil, err
	}
	value := result.InsertedID.(primitive.ObjectID)
	return &value, nil
}