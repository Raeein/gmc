package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Logger struct {
	Client     *mongo.Client
	DB         string
	Collection string
}

type Log struct {
	Title       string `bson:"title"`
	Description string `bson:"description"`
	Time        string `bson:"time"`
}

const timeout = 10 * time.Second

func New(username, password, db, collection string) Logger {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf(
		"mongodb+srv://%s:%s@gmc-logs.lp2jrjo.mongodb.net/?retryWrites=true&w=majority",
		username,
		password,
	)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return Logger{Client: client, DB: db, Collection: collection}
}

func (l Logger) Log(title, description string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	location, _ := time.LoadLocation("EST")
	currentTime := time.Now().In(location).Format("2006-01-02 15:04:05")
	// Check if bson works
	doc := Log{
		Title:       title,
		Description: description,
		Time:        currentTime,
	}
	collection := l.Client.Database(l.DB).Collection(l.Collection)
	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Inserted new document successfully", res.InsertedID)
}

func (l Logger) Close() {
	if err := l.Client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mongo connection closed.")
}
