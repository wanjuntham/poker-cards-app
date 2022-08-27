package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"testing"
)

func TestDatabaseOperations(t *testing.T) {
	_ = godotenv.Load("test.env")
	SetupDb()
	t.Run("Create deck", func(t *testing.T) {
		expected := uuid.NewString()
		deck := Deck{
			DeckId:    expected,
			Shuffled:  false,
			Remaining: 0,
			Cards:     nil,
		}
		actual, err := InsertDeck(deck)
		if err != nil {
			t.Errorf("Failed to insert to MongoDB: %v", err)
		}
		if actual != expected {
			t.Errorf("DeckID is different from _id inserted in mongodb, expected: %v, actual: %v", expected, actual)
		}
	})
	t.Run("Open deck", func(t *testing.T) {
		expected := uuid.NewString()
		SeedDb(expected)
		actual, err := OpenDeck(expected)
		if err != nil {
			t.Errorf("Failed to get data from mongodb: %v", err)
		}
		if actual.DeckId != expected {
			t.Errorf("DeckID from mongodb is different from what is expected, expected: %v, actual: %v", expected, actual.DeckId)
		}
		if actual.Remaining != 2 {
			t.Errorf("Remaining cards in deck doesn't match expected count, expected: %v, actual: %v", 2, actual.Remaining)
		}
	})
	t.Run("Open invalid deck", func(t *testing.T) {
		expected := uuid.NewString()
		SeedDb(uuid.NewString())
		actual, err := OpenDeck(expected)
		if err == nil {
			t.Errorf("An error is expected to return")
		}
		if actual.DeckId == expected {
			t.Errorf("Deck should not be found in DB, expected: %v, actual: %v", expected, actual.DeckId)
		}
	})
	t.Run("Open deck with empty DeckID", func(t *testing.T) {
		expected := ""
		SeedDb(uuid.NewString())
		actual, err := OpenDeck(expected)
		if err == nil {
			t.Errorf("An error is expected to return")
		}
		if actual.DeckId != expected {
			t.Errorf("Deck should not be found in DB, expected: %v, actual: %v", expected, actual.DeckId)
		}
	})
	t.Run("Draw 0 cards", func(t *testing.T) {
		expected := 0
		deckId := uuid.NewString()
		SeedDb(deckId)
		actual, err := DrawCardsFromDeck(deckId, expected)
		deck, err := OpenDeck(deckId)
		if err != nil {
			t.Errorf("Failed to get data from mongodb: %v", err)
		}
		if deck.Remaining != 2 {
			t.Errorf("Remaining cards in deck doesn't match expected count, expected: %v, actual: %v", 2, deck.Remaining)
		}
		if len(actual) != expected {
			t.Errorf("Cards drew is not matching, expected: %v, actual %v", expected, len(actual))
		}
	})
	t.Run("Draw 1 card", func(t *testing.T) {
		expected := 1
		deckId := uuid.NewString()
		SeedDb(deckId)
		actual, err := DrawCardsFromDeck(deckId, expected)
		deck, err := OpenDeck(deckId)
		if err != nil {
			t.Errorf("Failed to get data from mongodb: %v", err)
		}
		if deck.Remaining != 1 {
			t.Errorf("Remaining cards in deck doesn't match expected count, expected: %v, actual: %v", 1, deck.Remaining)
		}
		if len(actual) != expected {
			t.Errorf("Cards drew is not matching, expected: %v, actual %v", expected, len(actual))
		}
	})
	t.Run("Draw multiple cards", func(t *testing.T) {
		expected := 2
		deckId := uuid.NewString()
		SeedDb(deckId)
		actual, err := DrawCardsFromDeck(deckId, expected)
		deck, err := OpenDeck(deckId)
		if err != nil {
			t.Errorf("Failed to get data from mongodb: %v", err)
		}
		if deck.Remaining != 0 {
			t.Errorf("Remaining cards in deck doesn't match expected count, expected: %v, actual: %v", 0, deck.Remaining)
		}
		if len(actual) != expected {
			t.Errorf("Cards drew is not matching, expected: %v, actual %v", expected, len(actual))
		}
	})
	t.Run("Draw -1 card", func(t *testing.T) {
		expected := 0
		deckId := uuid.NewString()
		SeedDb(deckId)
		actual, err := DrawCardsFromDeck(deckId, -1)
		if err == nil {
			t.Errorf("An error is expected to return")
		}
		if len(actual) != expected {
			t.Errorf("Cards drew should be empty, expected: %v, actual %v", expected, len(actual))
		}
	})
	t.Run("Draw more cards than a deck has", func(t *testing.T) {
		expected := 0
		deckId := uuid.NewString()
		SeedDb(deckId)
		actual, err := DrawCardsFromDeck(deckId, 3)
		if err == nil {
			t.Errorf("An error is expected to return")
		}
		if len(actual) != expected {
			t.Errorf("Cards drew should be empty, expected: %v, actual %v", expected, len(actual))
		}
	})
	t.Run("Draw card from invalid deck", func(t *testing.T) {
		expected := 0
		deckId := uuid.NewString()
		SeedDb(uuid.NewString())
		actual, err := DrawCardsFromDeck(deckId, 3)
		if err == nil {
			t.Errorf("An error is expected to return")
		}
		if len(actual) != expected {
			t.Errorf("Cards drew should be empty, expected: %v, actual %v", expected, len(actual))
		}
	})
	t.Run("Draw card with empty DeckID", func(t *testing.T) {
		expected := 0
		deckId := ""
		SeedDb(uuid.NewString())
		actual, err := DrawCardsFromDeck(deckId, 3)
		if err == nil {
			t.Errorf("An error is expected to return")
		}
		if len(actual) != expected {
			t.Errorf("Cards drew should be empty, expected: %v, actual %v", expected, len(actual))
		}
	})
	TeardownDb()
}

func SeedDb(deckId string) {
	decks := []Deck{
		{
			DeckId:    deckId,
			Remaining: 2,
			Cards: []Card{
				{
					Value: "ACE",
					Suit:  "HEARTS",
					Code:  "AH",
				},
				{
					Value: "2",
					Suit:  "HEARTS",
					Code:  "2H",
				},
			},
		},
	}

	var allDecks []interface{}
	for _, d := range decks {
		allDecks = append(allDecks, d)
	}
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("POKER_COLLECTION_NAME")
	coll := DbClient.Database(dbName).Collection(collectionName)
	_, _ = coll.InsertMany(context.TODO(), allDecks)
}

func TeardownDb() {
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("POKER_COLLECTION_NAME")
	coll := DbClient.Database(dbName).Collection(collectionName)
	_, _ = coll.DeleteMany(context.TODO(), bson.D{})
	_ = DbClient.Disconnect(context.TODO())
	fmt.Println("closed connection to MongoDB")
}
