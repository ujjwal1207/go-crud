package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	json.NewDecoder(r.Body).Decode(&item)

	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(item)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var items []Item
	if err = cursor.All(ctx, &items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	objID, _ := primitive.ObjectIDFromHex(idParam)

	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item Item
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	objID, _ := primitive.ObjectIDFromHex(idParam)

	var item Item
	json.NewDecoder(r.Body).Decode(&item)

	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"name": item.Name}}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.ID = objID
	json.NewEncoder(w).Encode(item)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	objID, _ := primitive.ObjectIDFromHex(idParam)

	collection := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Deleted successfully"))
}
