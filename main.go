package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/iwayankurniawan/gocms/models"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func readContent(w http.ResponseWriter, r *http.Request) {
	godotenv.Load(".env")
	// create an aws session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-north-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("aws_access_key_id"), os.Getenv("aws_secret_access_key"), ""),
	},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	params := mux.Vars(r)
	id := params["id"]
	item := models.Content{}

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("content"), //table name
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		response := models.Response{Status: http.StatusInternalServerError, Message: "Got error calling GetItem", Data: map[string]interface{}{"data": err.Error()}}
		json.NewEncoder(w).Encode(response)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	if item.Text == "" {
		fmt.Println("Could not find data")
		w.WriteHeader(http.StatusOK)
		response := models.Response{Status: http.StatusInternalServerError, Message: "Could not find data", Data: map[string]interface{}{"data": ""}}
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Println("Found item:")
	fmt.Println("Id:  ", item.Id)
	fmt.Println("Text: ", item.Text)

	w.WriteHeader(http.StatusOK)
	response := models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": item}}
	json.NewEncoder(w).Encode(response)
}

func createContent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	content := models.Content{Id: id, Text: "Hello Yorobun"}

	w.WriteHeader(http.StatusOK)
	response := models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": content}}
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/project/{id}", readContent).Methods("GET")

	http.ListenAndServe(":8000", r)
}
