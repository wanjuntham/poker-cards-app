package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/decks", func(context *gin.Context) {
		var shuffled bool
		var requestedCards []string
		var err error
		shuffleString, exists := context.GetQuery("shuffle")
		if !exists {
			shuffled = false
		} else {
			if shuffled, err = strconv.ParseBool(shuffleString); err != nil {
				context.JSON(http.StatusBadRequest, err.Error())
				return
			}
		}
		requestedCardsString, exists := context.GetQuery("cards")
		if exists {
			requestedCards = strings.Split(requestedCardsString, ",")
		}

		result, err := CreateDeck(shuffled, requestedCards...)
		if err != nil {
			context.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		context.JSON(http.StatusCreated, result)
	})

	r.GET("/decks/:deckId", func(context *gin.Context) {
		deckId := context.Param("deckId")
		_, err := uuid.Parse(deckId)
		if err != nil {
			context.JSON(http.StatusBadRequest, err.Error())
			return
		}

		result, err := OpenDeck(deckId)
		if err != nil {
			context.JSON(http.StatusNotFound, err.Error())
			return
		}
		context.JSON(http.StatusOK, result)
	})

	r.GET("/decks/:deckId/cards/count/:count", func(context *gin.Context) {
		deckId := context.Param("deckId")
		_, err := uuid.Parse(deckId)
		if err != nil {
			context.JSON(http.StatusBadRequest, err.Error())
			return
		}

		count, err := strconv.Atoi(context.Param("count"))
		if err != nil {
			context.JSON(http.StatusBadRequest, err.Error())
			return
		}

		result, err := DrawCards(deckId, count)
		if err != nil {
			context.JSON(http.StatusNotFound, err.Error())
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"cards": result,
		})
	})

	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := SetupRouter()
	SetupDb()
	defer func() {
		if err := DbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
