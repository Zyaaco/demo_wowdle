package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const singletonID = "current_word"

type word_of_the_day struct {
	ID        string    `bson:"_id"`
	Word      string    `bson:"word"`
	UpdatedAt time.Time `bson:"updated_at"`
}

const collName = "word_of_the_day"

func ChangeWord(newWord string) error {
	collection := mongoClient.Database(db).Collection(collName)

	// look for specific ID
	filter := bson.M{"_id": singletonID}

	update := bson.M{"$set": bson.M{
		"word":       newWord,
		"updated_at": time.Now(),
	}}

	opts := options.UpdateOne().SetUpsert(true)

	// Execute
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func GetWord() (string, error) {
	collection := mongoClient.Database(db).Collection(collName)

	filter := bson.M{"_id": singletonID}

	var result word_of_the_day

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Word, nil
}
