# Poker Game API

An API to handle poker deck and cards.

## Prerequisites

### Go

[Please install go from here.](https://go.dev/dl/)

### MongoDB

There are multiple ways to run a MongoDB instance. The preferred way will be to run with Docker.

If you have Docker Desktop installed and started, run below command in your terminal on your root directory.
```shell
docker compose up -d
```

MongoDB is ready when you see this in your terminal.
```
[+] Running 2/2
 - Network golang-assignment_net        Created                                                                                                                                                                                                                                                                          0.0s 
 - Container golang-assignment-mongo-1  Started
```

You can shut down MongoDB by running command below.
```shell
docker compose down
```

### Docker

_Required only if you want to run MongoDB on docker._

## Install dependencies (Optional)

Dependencies should be installed automatically during `go run` / `go build`.
 ```shell
 go mod download
 ```

## Quick start

There are 3 ways to run the application:

### Go run

```shell
go run .
```

### Go build and execute

```shell
go build

# if you are using windows
# Command Prompt
card-game.exe
# powershell
.\card-game.exe

# if you are using linux
./card-game
```

### Go install

```shell
go install

card-game
```

## Usage

Server is exposed on port 8080, please make HTTP requests to this base URL: http://localhost:8080

| Use               | Relative endpoint                     | Local absolute endpoint                                  | HTTP Method | Query supported                                     |
|-------------------|---------------------------------------|----------------------------------------------------------|-------------|-----------------------------------------------------|
| Create a new deck | `/decks`                              | http://localhost:8080/decks                              | POST        | `shuffle`: `true`/`false`<br/>`cards`: `AD`/`AD,KH` |
| Open a deck       | `/decks/{DeckID}`                     | http://localhost:8080/decks/{deckID}                     | GET         | N/A                                                 |
| Draw a card       | `/decks/{DeckID}/cards/count/{count}` | http://localhost:8080/decks/{deckID}/cards/count/{count} | GET         | N/A                                                 |

*Please substitute `{DeckID}` with the `deck_id` you received in the `Create a new deck`'s response body.

*Please substitute `{count}` with how many cards you want to draw from the deck.

## Test

### Preparation

- [MongoDB server running locally](#MongoDB)

### Command

```shell
# To run all tests
go test

# If you want to view the status for each job
go test -v
```
