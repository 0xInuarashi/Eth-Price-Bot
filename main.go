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
			fmt.Printf("WorkerId %v got %v \n", id, "$"+res)
			err = discord.UpdateWatchStatus(0, "fast:"+res.Data.Fast+"medium:"+res.Data.Medium+"slow:"+res.Data.Slow)
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
	Data struct {
		Fast     string `json:"fastestFee"`
		Medium string `json:"halfHourFee"`
		Slow   string `json:"hourFee"`
	} `json:"data"`
}

func getPrice() (Response, error) {
	res, err := http.Get(endpoint)

	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return "", fmt.Errorf("failed to fetch: %v", err)
	}

	jsonPayload, err := decodeJson[Response](res.Body)

	fmt.Println(jsonPayload)

	if err != nil {
		return "", fmt.Errorf("failed to decode json: %v", err)
	}

	fast, err := strconv.ParseFloat(jsonPayload.Data.Fast, 64)
	if err != nil {
		return "", fmt.Errorf("invalid amount format: %v", err)
	}

	medium, err := strconv.ParseFloat(jsonPayload.Data.Medium, 64)
	if err != nil {
		return "", fmt.Errorf("invalid amount format: %v", err)
	}

	slow, err := strconv.ParseFloat(jsonPayload.Data.Slow, 64)
	if err != nil {
		return "", fmt.Errorf("invalid amount format: %v", err)
	}

	return Response{
		Fast: fmt.Sprintf("%.0f", fast),
		Medium: fmt.Sprintf("%.0f", medium),
		Slow: fmt.Sprintf("%.0f", slow),
	}
}

func decodeJson[T any](r io.Reader) (T, error) {
	var v T
	err := json.NewDecoder(r).Decode(&v)
	return v, err
}
