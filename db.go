package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var DbClient *mongo.Client

func SetupDb() {
	connectionString := os.Getenv("MONGO_CONNECTION_STRING")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		panic(err)
	}
	DbClient = client
}

func InsertDeck(deck Deck) (interface{}, error) {
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("POKER_COLLECTION_NAME")
	coll := DbClient.Database(dbName).Collection(collectionName)
	result, err := coll.InsertOne(context.TODO(), deck)
	if err != nil {
		return result, err
	}
	return result.InsertedID, err
}

func GetDeck(deckId string) (Deck, error) {
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("POKER_COLLECTION_NAME")
	coll := DbClient.Database(dbName).Collection(collectionName)
	var result Deck
	err := coll.FindOne(context.TODO(), bson.D{{"_id", deckId}}).Decode(&result)
	return result, err
}

func DrawCardsFromDeck(deckId string, count int) ([]Card, error) {
	var cards []Card

	deck, err := GetDeck(deckId)
	if err != nil {
		return nil, err
	}

	if len(deck.Cards) < count {
		return nil, errors.New("deck has less cards than count intended to draw")
	}

	for i := 0; i < count; i++ {
		var removedCard Card
		removedCard, deck.Cards = deck.Cards[len(deck.Cards)-1], deck.Cards[:len(deck.Cards)-1]
		cards = append(cards, removedCard)
	}
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("POKER_COLLECTION_NAME")
	coll := DbClient.Database(dbName).Collection(collectionName)

	var cardCodes []string
	for _, card := range cards {
		cardCodes = append(cardCodes, card.Code)
	}

	update := bson.D{{
		"$pull",
		bson.D{{
			"cards",
			bson.D{{
				"code",
				bson.D{{
					"$in",
					cardCodes,
				}},
			}},
		}},
	}, {
		"$set",
		bson.D{{
			"remaining",
			len(deck.Cards),
		}},
	}}

	_, err = coll.UpdateOne(context.TODO(), bson.D{{"_id", deckId}}, update)

	if err != nil {
		return nil, err
	}
	return cards, nil
}
