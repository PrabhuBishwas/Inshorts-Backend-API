package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
  "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

type Article struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Subtitle  string             `json:"subtitle,omitempty" bson:"subtitle,omitempty"`
  Content   string             `json:"content,omitempty" bson:"content,omitempty"`
  Date  time.Time
}

func CreateArticleEndpoint(response http.ResponseWriter, request *http.Request) {
  switch request.Method {
  case "GET":
    response.Header().Set("content-type", "application/json")
  	var articles []Article
  	collection := client.Database("inshort").Collection("article")
  	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
  	cursor, err := collection.Find(ctx, bson.M{})
  	if err != nil {
  		response.WriteHeader(http.StatusInternalServerError)
  		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
  		return
  	}
  	defer cursor.Close(ctx)
  	for cursor.Next(ctx) {
  		var article Article
  		cursor.Decode(&article)
  		articles = append(articles, article)
  	}
  	if err := cursor.Err(); err != nil {
  		response.WriteHeader(http.StatusInternalServerError)
  		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
  		return
  	}
  	json.NewEncoder(response).Encode(articles)

  case "POST":
    response.Header().Set("content-type", "application/json")
  	var article Article
  	_ = json.NewDecoder(request.Body).Decode(&article)
  	collection := client.Database("inshort").Collection("article")
  	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    article.Date = time.Now()
  	result, _ := collection.InsertOne(ctx, article)
  	json.NewEncoder(response).Encode(result)
  }
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	http.HandleFunc("/article", CreateArticleEndpoint)
  err := http.ListenAndServe(":9090", nil) // set listen port
  if err != nil {
      log.Fatal("ListenAndServe: ", err)
  }
}
