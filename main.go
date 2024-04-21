package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var endpoint = "https://mempool.space/api/v1/fees/recommended" // unused lol
const shards = 1

func worker(id int, token string, coin string) {
	discord, err := discordgo.New("Bot " + token)
	endpoint = "https://mempool.space/api/v1/fees/recommended" 

	if err != nil {
		log.Fatalf("Error creating discord session: %v", err)
	}

	discord.ShardCount = shards
	discord.ShardID = id

	err = discord.Open()
	if err != nil {
		log.Fatalf("Error opening discord ws: %v", err)
	}
	defer discord.Close()

	for {
		res, err := getPrice()
		if err != nil {
			log.Printf("Error getting price for shard %d: %v \n", id, err)
		} else {
			// fmt.Printf("WorkerId %v got %v \n", id, "$"+res)
			err = discord.UpdateWatchStatus(0, "⚡"+strconv.Itoa(res.Fast)+" |🚶‍♂️"+strconv.Itoa(res.Medium)+" |🐢"+strconv.Itoa(res.Slow))
			if err != nil {
				log.Printf("Error updating discord status for shard %d: %v \n", id, err)
			}
		}
		time.Sleep(30 * time.Second)
	}

}

func main() {
	fmt.Println("hello world 🌍👋")
	token := getEnvOrDie("TOKEN")
	coin := getEnvOrDie("COIN")

	wg := sync.WaitGroup{}

	for shardId := 0; shardId < shards; shardId++ {
		wg.Add(1)
		go worker(shardId, token, coin)
	}

	wg.Wait()
}

func getEnvOrDie(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Could not find %v in .env", key)
	}

	return value
}

type Response struct {
	Fast   int `json:"fastestFee"`
	Medium int `json:"halfHourFee"`
	Slow   int `json:"hourFee"`
}

type Return struct {
	Fast   string `json:"fastestFee"`
	Medium string `json:"halfHourFee"`
	Slow   string `json:"hourFee"`
}

func getPrice() (Response, error) {
	res, err := http.Get(endpoint)

	if res != nil {
		defer res.Body.Close()
	}

	fmt.Println(res)
	fmt.Println("Body")
	fmt.Println(res.Body)

	// if err != nil {
	// 	return "", fmt.Errorf("failed to fetch: %v", err)
	// }

	// jsonPayload, err := decodeJson[Response](res.Body)
	var jsonPayload Response
	err = json.NewDecoder(res.Body).Decode(&jsonPayload)

	if err != nil {
        fmt.Printf("Error decoding JSON: %s\n", err)
        return Response{}, fmt.Errorf("fail")
    }

    fmt.Printf("Decoded JSON: %+v\n", jsonPayload)

	fmt.Println("Hi")
	fmt.Println(jsonPayload)

	if err != nil {
		return Response{}, fmt.Errorf("failed to decode json: %v", err)
	}

	// fast, err := strconv.ParseFloat(jsonPayload.Fast, 64)
	fast := jsonPayload.Fast
	if err != nil {
		return Response{}, fmt.Errorf("invalid amount format: %v", err)
	}

	// medium, err := strconv.ParseFloat(jsonPayload.Medium, 64)
	medium := jsonPayload.Medium
	if err != nil {
		return Response{}, fmt.Errorf("invalid amount format: %v", err)
	}

	// slow, err := strconv.ParseFloat(jsonPayload.Slow, 64)
	slow := jsonPayload.Slow
	if err != nil {
		return Response{}, fmt.Errorf("invalid amount format: %v", err)
	}

	return Response{
		Fast:   fast,
		Medium: medium,
		Slow:   slow,
	}, nil	
}

func decodeJson[T any](r io.Reader) (T, error) {
	var v T
	err := json.NewDecoder(r).Decode(&v)
	return v, err
}
