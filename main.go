package main

import (
	"context"
	"encoding/json"

	// "reflect"

	// "io/ioutil"
	"net/http"

	// "encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type Post struct {
	ID       primitive.ObjectID `bson:"_id"`
	Caption  string             `bson:"caption"`
	ImageURL string             `bson:"imageURL"`
	PostedAt string             `bson:"postedAt"`
	PostedBy string             `bson:"postedBy"`
}

type Message struct {
	Success bool   `bson:"success"`
	Message string `bson:"message,omitbody"`
}

func (u *User) fillUser() {
	u.ID = primitive.NewObjectID()
}

func (p *Post) fillPost() {
	p.ID = primitive.NewObjectID()
	p.PostedAt = time.Now().String()
}

func ConnectDB() (*mongo.Collection, *mongo.Collection) {

	// Set client options
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://trialuser:qaz12345@cluster0.hlrin.mongodb.net/appointy1?retryWrites=true&w=majority"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongoDb")

	appointyDatabase := client.Database("appointy1")
	userCollection := appointyDatabase.Collection("users")
	postCollection := appointyDatabase.Collection("posts")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return userCollection, postCollection
}

var userCollection, postCollection = ConnectDB()

func createUser(res http.ResponseWriter, req *http.Request) {
	var user User
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	if len(user.Password) > 0 && len(user.Email) > 0 && len(user.Name) > 0 {
		user.fillUser()
		var msg Message

		userResult, err := userCollection.InsertOne(ctx, bson.D{
			{"_id", user.ID},
			{"name", user.Name},
			{"email", user.Email},
			{"password", user.Password},
		})
		if err != nil {
			log.Fatal(err)
		}
		msg.Success = true
		msg.Message = "The user was succesfully added"
		fmt.Printf("%+v", userResult)
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(msg)
	} else {
		var msg Message
		msg.Success = false
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(msg)
	}
}

func createPost(res http.ResponseWriter, req *http.Request) {

	var post Post
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := json.NewDecoder(req.Body).Decode(&post)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(res, "Post: %+v", post)

	if len(post.Caption) > 0 && len(post.ImageURL) > 0 {
		post.fillPost()

		postResult, err := postCollection.InsertOne(ctx, bson.D{
			{"_id", post.ID},
			{"caption", post.Caption},
			{"imageURL", post.ImageURL},
			{"postedAt", post.PostedAt},
			{"postedBy", post.PostedBy},
		})
		if err != nil {
			log.Fatal(err)
		}
		var msg Message
		msg.Success = true
		msg.Message = "The post was succesfully created"
		fmt.Printf("%+v", postResult)
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(msg)
	} else {
		var msg Message
		msg.Success = false
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(msg)
	}
}

func getPost(res http.ResponseWriter, req *http.Request) {
	postId := req.URL.Path[len("/posts/"):]
	objectId, err := primitive.ObjectIDFromHex(postId)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := postCollection.Find(ctx, bson.M{"_id": objectId})
	if err != nil {
		log.Fatal(err)
	}
	var post []bson.M
	if err = cursor.All(ctx, &post); err != nil {
		log.Fatal(err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(post)
}

func getUser(res http.ResponseWriter, req *http.Request) {
	userId := req.URL.Path[len("/users/"):]
	objectId, err := primitive.ObjectIDFromHex(userId)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := userCollection.Find(ctx, bson.M{"_id": objectId})
	if err != nil {
		log.Fatal(err)
	}
	var user []bson.M
	if err = cursor.All(ctx, &user); err != nil {
		log.Fatal(err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(user)
}

func getPostsByUser(res http.ResponseWriter, req *http.Request) {
	userId := req.URL.Path[len("/posts/users/"):]
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := postCollection.Find(ctx, bson.M{"postedBy": userId})
	if err != nil {
		log.Fatal(err)
	}
	var posts []bson.M
	if err = cursor.All(ctx, &posts); err != nil {
		log.Fatal(err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(posts)
}

func main() {

	http.HandleFunc("/users", createUser)
	http.HandleFunc("/posts", createPost)
	http.HandleFunc("/posts/", getPost)
	http.HandleFunc("/users/", getUser)
	http.HandleFunc("/posts/users/", getPostsByUser)
	http.ListenAndServe(":3000", nil)

	// cursor, err := userCollection.Find(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var users []bson.D
	// if err = cursor.All(ctx, &users); err != nil {
	// 	log.Fatal(err)
	// }
	// newUser := bson.Unmarshal(users)
	// jsonData, err := json.MarshalIndent(users, "", "    ")
	// fmt.Println(users)

}
