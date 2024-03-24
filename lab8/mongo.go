// Example use of Go mongo-driver
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbEndpoint = "mongodb://172.17.0.2:27017" // Find this from the Mongo container
)

type Item struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Price float64            `bson:"price"`
}

var (
	client     *mongo.Client
	collection *mongo.Collection
)

func main() {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Select database and collection
	collection = client.Database("store").Collection("items")

	// Initialize router
	router := http.NewServeMux()

	// Define routes
	router.HandleFunc("/list", listItems)
	router.HandleFunc("/price", getItemPrice)
	router.HandleFunc("/create", createItem)
	router.HandleFunc("/update", updateItem)
	router.HandleFunc("/delete", deleteItem)

	// Start server
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func listItems(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var items []Item
	if err := cursor.All(ctx, &items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, item := range items {
		fmt.Fprintf(w, "%s: $%.2f\n", item.Name, item.Price)
	}
}

func getItemPrice(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name := r.URL.Query().Get("name")

	var item Item
	if err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&item); err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Price of %s: $%.2f\n", item.Name, item.Price)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name := r.URL.Query().Get("name")
	priceStr := r.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	item := Item{Name: name, Price: price}
	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Item created with ID %s\n", result.InsertedID.(primitive.ObjectID).Hex())
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name := r.URL.Query().Get("name")
	priceStr := r.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	update := bson.M{"$set": bson.M{"price": price}}
	result, err := collection.UpdateOne(ctx, bson.M{"name": name}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Updated %d item(s)\n", result.ModifiedCount)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name := r.URL.Query().Get("name")
	result, err := collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Deleted %d item(s)\n", result.DeletedCount)
}
