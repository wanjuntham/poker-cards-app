package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	_ = godotenv.Load("test.env")
	router := SetupRouter()
	SetupDb()
	t.Run("Create deck", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/decks", nil)

		router.ServeHTTP(w, req)

		var resBody Deck
		if resBodyBytes := w.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		if resBody.Shuffled {
			t.Errorf("Deck should not be shuffled. Expected: %v, actual: %v", false, resBody.Shuffled)
		}

		if resBody.Remaining != 52 {
			t.Errorf("Deck remaining cards should be 52. Expected: %v, actual: %v", 52, resBody.Remaining)
		}

		if _, err := uuid.Parse(resBody.DeckId); err != nil {
			t.Errorf("Deck ID should be a UUID. Expected: random UUID, actual: %v", resBody.DeckId)
		}

		if w.Code != http.StatusCreated {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusCreated, w.Code)
		}
	})
	t.Run("Create shuffled deck", func(t *testing.T) {

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		q := req.URL.Query()
		q.Add("shuffle", "true")
		req.URL.RawQuery = q.Encode()

		router.ServeHTTP(w, req)

		var resBody Deck
		if resBodyBytes := w.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		if !resBody.Shuffled {
			t.Errorf("Deck should be shuffled. Expected: %v, actual: %v", true, resBody.Shuffled)
		}

		if resBody.Remaining != 52 {
			t.Errorf("Deck remaining cards should be 52. Expected: %v, actual: %v", 52, resBody.Remaining)
		}

		if _, err := uuid.Parse(resBody.DeckId); err != nil {
			t.Errorf("Deck ID should be a UUID. Expected: random UUID, actual: %v", resBody.DeckId)
		}

		if w.Code != http.StatusCreated {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusCreated, w.Code)
		}
	})
	t.Run("Create partial deck", func(t *testing.T) {
		partialCardsStr := "AH,2D"
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		q := req.URL.Query()
		q.Add("cards", partialCardsStr)
		req.URL.RawQuery = q.Encode()

		router.ServeHTTP(w, req)

		var resBody Deck
		if resBodyBytes := w.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		if resBody.Shuffled {
			t.Errorf("Deck should not be shuffled. Expected: %v, actual: %v", false, resBody.Shuffled)
		}

		if resBody.Remaining != 2 {
			t.Errorf("Deck remaining cards should be 2. Expected: %v, actual: %v", 2, resBody.Remaining)
		}

		if _, err := uuid.Parse(resBody.DeckId); err != nil {
			t.Errorf("Deck ID should be a UUID. Expected: random UUID, actual: %v", resBody.DeckId)
		}

		if w.Code != http.StatusCreated {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusCreated, w.Code)
		}
	})
	t.Run("Create shuffled deck with invalid query string", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		q := req.URL.Query()
		q.Add("shuffle", "random value")
		req.URL.RawQuery = q.Encode()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusBadRequest, w.Code)
		}
	})
	t.Run("Create partial deck with invalid query string", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		q := req.URL.Query()
		q.Add("cards", "random value")
		req.URL.RawQuery = q.Encode()

		router.ServeHTTP(w, req)

		var resBody Deck
		if resBodyBytes := w.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		if w.Code != http.StatusCreated {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusCreated, w.Code)
		}

		if resBody.Remaining != 0 {
			t.Errorf("Deck remaining should be 0. Expected: %v, actual: %v", 0, resBody.Remaining)
		}
	})
	t.Run("Open deck", func(t *testing.T) {
		//arrange
		seedW := httptest.NewRecorder()
		deckReq, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		router.ServeHTTP(seedW, deckReq)
		var seedBody Deck
		if resBodyBytes := seedW.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &seedBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		//act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s", seedBody.DeckId), nil)

		router.ServeHTTP(w, req)

		var resBody Deck
		if resBodyBytes := w.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		if resBody.Shuffled {
			t.Errorf("Deck should not be shuffled. Expected: %v, actual: %v", false, resBody.Shuffled)
		}

		if resBody.Remaining != 52 {
			t.Errorf("Deck remaining cards should be 52. Expected: %v, actual: %v", 52, resBody.Remaining)
		}

		if _, err := uuid.Parse(resBody.DeckId); err != nil {
			t.Errorf("Deck ID should be a UUID. Expected: random UUID, actual: %v", resBody.DeckId)
		}

		if w.Code != http.StatusOK {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusOK, w.Code)
		}

		if len(resBody.Cards) != 52 {
			t.Errorf("Deck default cards count should be 52. Expected: %v, actual: %v", 52, len(resBody.Cards))
		}
	})
	t.Run("Open invalid deck", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s", uuid.NewString()), nil)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusNotFound, w.Code)
		}
	})
	t.Run("Open shuffled deck", func(t *testing.T) {
		//arrange
		seedW := httptest.NewRecorder()
		deckReq, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		q := deckReq.URL.Query()
		q.Add("shuffle", "true")
		deckReq.URL.RawQuery = q.Encode()
		router.ServeHTTP(seedW, deckReq)
		var seedBody Deck
		if resBodyBytes := seedW.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &seedBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		//act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s", seedBody.DeckId), nil)

		router.ServeHTTP(w, req)

		var resBody Deck
		if resBodyBytes := w.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		if !resBody.Shuffled {
			t.Errorf("Deck should be shuffled. Expected: %v, actual: %v", true, resBody.Shuffled)
		}

		if resBody.Remaining != 52 {
			t.Errorf("Deck remaining cards should be 52. Expected: %v, actual: %v", 52, resBody.Remaining)
		}

		if _, err := uuid.Parse(resBody.DeckId); err != nil {
			t.Errorf("Deck ID should be a UUID. Expected: random UUID, actual: %v", resBody.DeckId)
		}

		if w.Code != http.StatusOK {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusOK, w.Code)
		}

		if len(resBody.Cards) != 52 {
			t.Errorf("Deck default cards count should be 52. Expected: %v, actual: %v", 52, len(resBody.Cards))
		}

		if !isShuffled(resBody.Cards) {
			t.Error("Deck should be shuffled, please check the sequence of the cards.")
		}
	})
	t.Run("Open partial deck", func(t *testing.T) {
		//arrange
		partialCardsStr := "AH,2D"
		seedW := httptest.NewRecorder()
		deckReq, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		q := deckReq.URL.Query()
		q.Add("cards", partialCardsStr)
		deckReq.URL.RawQuery = q.Encode()
		router.ServeHTTP(seedW, deckReq)
		var seedBody Deck
		if resBodyBytes := seedW.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &seedBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		//act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s", seedBody.DeckId), nil)

		router.ServeHTTP(w, req)

		var resBody Deck
		if resBodyBytes := w.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		if resBody.Shuffled {
			t.Errorf("Deck should not be shuffled. Expected: %v, actual: %v", true, resBody.Shuffled)
		}

		if resBody.Remaining != 2 {
			t.Errorf("Deck remaining cards should be 52. Expected: %v, actual: %v", 2, resBody.Remaining)
		}

		if _, err := uuid.Parse(resBody.DeckId); err != nil {
			t.Errorf("Deck ID should be a UUID. Expected: random UUID, actual: %v", resBody.DeckId)
		}

		if w.Code != http.StatusOK {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusOK, w.Code)
		}

		if len(resBody.Cards) != 2 {
			t.Errorf("Deck default cards count should be 52. Expected: %v, actual: %v", 2, len(resBody.Cards))
		}

		partialCards := strings.Split(partialCardsStr, ",")
		isMatch := true
		for _, c := range resBody.Cards {
			for i, code := range partialCards {
				if c.Code == code {
					break
				}
				if i == len(partialCards)-1 {
					isMatch = false
				}
			}
		}
		if !isMatch {
			t.Error("Requested partial cards should match cards in deck")
		}
	})
	t.Run("Draw 0 cards", func(t *testing.T) {
		//arrange
		seedW := httptest.NewRecorder()
		deckReq, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		router.ServeHTTP(seedW, deckReq)
		var seedBody Deck
		if resBodyBytes := seedW.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &seedBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		//act
		count := 0
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s/cards/count/%d", seedBody.DeckId, count), nil)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusNotFound, w.Code)
		}
	})
	t.Run("Draw 1 card", func(t *testing.T) {
		//arrange
		seedW := httptest.NewRecorder()
		deckReq, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		router.ServeHTTP(seedW, deckReq)
		var seedBody Deck
		if resBodyBytes := seedW.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &seedBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		//act
		count := 1
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s/cards/count/%d", seedBody.DeckId, count), nil)

		router.ServeHTTP(w, req)

		var resBody Deck
		if resBodyBytes := w.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		if w.Code != http.StatusOK {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusOK, w.Code)
		}

		if len(resBody.Cards) != count {
			t.Errorf("Length of cards is incorrect. expected: %v, actual %v", count, len(resBody.Cards))
		}
	})
	t.Run("Draw multiple cards", func(t *testing.T) {
		//arrange
		seedW := httptest.NewRecorder()
		deckReq, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		router.ServeHTTP(seedW, deckReq)
		var seedBody Deck
		if resBodyBytes := seedW.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &seedBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		//act
		count := 5
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s/cards/count/%d", seedBody.DeckId, count), nil)

		router.ServeHTTP(w, req)

		var resBody Deck
		if resBodyBytes := w.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		if w.Code != http.StatusOK {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusOK, w.Code)
		}

		if len(resBody.Cards) != count {
			t.Errorf("Length of cards is incorrect. expected: %v, actual %v", count, len(resBody.Cards))
		}
	})
	t.Run("Draw -1 card", func(t *testing.T) {
		//arrange
		seedW := httptest.NewRecorder()
		deckReq, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		router.ServeHTTP(seedW, deckReq)
		var seedBody Deck
		if resBodyBytes := seedW.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &seedBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		//act
		count := -1
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s/cards/count/%d", seedBody.DeckId, count), nil)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusNotFound, w.Code)
		}
	})
	t.Run("Draw more cards than a deck has", func(t *testing.T) {
		//arrange
		seedW := httptest.NewRecorder()
		deckReq, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		router.ServeHTTP(seedW, deckReq)
		var seedBody Deck
		if resBodyBytes := seedW.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &seedBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		//act
		count := 53
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s/cards/count/%d", seedBody.DeckId, count), nil)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusOK, w.Code)
		}

	})
	t.Run("Draw cards from invalid deck", func(t *testing.T) {
		//arrange
		seedW := httptest.NewRecorder()
		deckReq, _ := http.NewRequest(http.MethodPost, "/decks", nil)
		router.ServeHTTP(seedW, deckReq)
		var seedBody Deck
		if resBodyBytes := seedW.Body.Bytes(); resBodyBytes != nil {
			if err := json.Unmarshal(resBodyBytes, &seedBody); err != nil {
				t.Error("Error while unmarshaling response body to Deck struct.")
			}
		}

		//act
		count := 1
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/decks/%s/cards/count/%d", uuid.NewString(), count), nil)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("HTTP status code is incorrect. expected: %v, actual: %v", http.StatusOK, w.Code)
		}

	})
	TeardownDb()
}

func isShuffled(cards []Card) bool {
	allValues := []string{"ACE", "2", "3", "4", "5", "6", "7", "8", "9", "10", "JACK", "QUEEN", "KING"}
	allSuits := []string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}
	var sequentialCards []Card

	for _, suit := range allSuits {
		for _, value := range allValues {
			firstCharValueRune := []rune(value)
			firstCharValue := string(firstCharValueRune[0:1])
			firstCharSuitRune := []rune(suit)
			firstCharSuit := string(firstCharSuitRune[0:1])
			cardCode := firstCharValue + firstCharSuit
			sequentialCards = append(sequentialCards, Card{
				Value: value,
				Suit:  suit,
				Code:  cardCode,
			})
		}
	}

	for i, _ := range sequentialCards {
		if cards[i].Code != sequentialCards[i].Code {
			return true
		}
	}
	return false
}
