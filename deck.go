package main

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Deck struct {
	DeckId    string `json:"deck_id" bson:"_id"`
	Shuffled  bool   `json:"shuffled" bson:"shuffled"`
	Remaining int    `json:"remaining" bson:"remaining"`
	Cards     []Card `json:"cards,omitempty" bson:"cards,omitempty"`
}

func (d *Deck) GenerateCards(shuffle bool, requestedCards ...string) {

	allValues := []string{"ACE", "2", "3", "4", "5", "6", "7", "8", "9", "10", "JACK", "QUEEN", "KING"}
	allSuits := []string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}
	var allCards []Card

	for _, suit := range allSuits {
		for _, value := range allValues {
			firstCharValueRune := []rune(value)
			firstCharValue := string(firstCharValueRune[0:1])
			firstCharSuitRune := []rune(suit)
			firstCharSuit := string(firstCharSuitRune[0:1])
			cardCode := firstCharValue + firstCharSuit
			allCards = append(allCards, Card{
				Value: value,
				Suit:  suit,
				Code:  cardCode,
			})
		}
	}

	// shuffle the card slice if shuffle is true
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(allCards), func(i, j int) {
			allCards[i], allCards[j] = allCards[j], allCards[i]
		})
	}

	// returns only requested cards
	if len(requestedCards) > 0 {
		var partialCards []Card
		for _, requestedCard := range requestedCards {
			for i, card := range allCards {
				if requestedCard == card.Code {
					partialCards = append(partialCards, allCards[i])
					break
				}
			}
		}

		d.Cards = partialCards
	} else {
		d.Cards = allCards
	}
}

func CreateDeck(shuffled bool, requestedCards ...string) (Deck, error) {
	deck := Deck{DeckId: uuid.NewString(), Shuffled: shuffled}
	deck.GenerateCards(shuffled, requestedCards...)
	deck.Remaining = len(deck.Cards)
	_, err := InsertDeck(deck)
	deck.Cards = nil
	if err != nil {
		return deck, err
	}
	return deck, nil
}

func OpenDeck(deckId string) (Deck, error) {
	od, err := GetDeck(deckId)
	if err != nil {
		return od, err
	}
	return od, nil
}

func DrawCards(deckId string, count int) ([]Card, error) {
	cards, err := DrawCardsFromDeck(deckId, count)
	if err != nil {
		return nil, err
	}
	return cards, nil
}
