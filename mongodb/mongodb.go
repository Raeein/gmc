package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Logger struct {
	Client *mongo.Client
	Ctx    context.Context
}

func New(username, password string) Logger {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("mongodb+srv://%s:%s@gmc-logs.lp2jrjo.mongodb.net/?retryWrites=true&w=majority", username, password)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return Logger{Client: client, Ctx: ctx}
}

func (l Logger) Log(db, cl string) {
	collection := l.Client.Database(db).Collection(cl)
	collection.InsertOne(l.Ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
}
