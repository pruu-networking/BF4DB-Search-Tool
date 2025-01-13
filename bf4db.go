package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var IsDebug = false
var homeDir, _ = os.UserHomeDir()
var envPath = filepath.Join(homeDir, "/.BF4DB-Search-Tool")
var apiKey string
var IsDiscord = false

func main() {
	// Verify if API key is set
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	key, err := godotenv.Read(envPath)
	if err != nil {
		fmt.Println("No API key found.")
		setApiKey()
	}
	apiKey = key["BF4DB_API_KEY"] // Set API key

	if len(os.Args) < 2 {
		fmt.Println("Usage: bf4db <player name> | bf4db <ip> to search your own ip | -h for help")
		fmt.Println("Use bf4db `<discord_id> dc` to search by Discord User ID")
		return
	}
	player := os.Args[1]

	if len(os.Args) > 2 {
		if os.Args[2] == "dbg" {
			IsDebug = true
		}
		if os.Args[2] == "dc" || os.Args[2] == "discord" || os.Args[2] == "-dc" {
			IsDiscord = true
		}

	}
	if player == "-config" || player == "-c" {
		// Set a new API key
		setApiKey()
	}
	if player == "-help" || player == "-h" || player == "-?" || player == "?" {
		// Show help
		fmt.Println("Usage: bf4db <player name> to search for a player or bf4db <ip> to search your own ip")
		fmt.Println("Use bf4db `<discord_id> dc` to search by Discord User ID")
		fmt.Println("Usage: bf4db -c to set a new API key")
		fmt.Println("Add 'dbg' or 'debug' to enable debug mode, not available for discord search")
		fmt.Println("Usage: bf4db -h to show this weird help message")
		return
	}
	// Check if is an ip:port, useful when CTRC+C players IP from Procon Layer
	ip, _, err := net.SplitHostPort(player)
	if err == nil {
		player = ip
	}

	// Search users own IP
	if player == "ip" || player == "IP" {
		resp, err := http.Get("https://api.ipify.org")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		player = string(body)
	}
	fmt.Println("Searching for " + player + "\n")
	GlobalSearch(player)
}

func setApiKey() {

	// Prompt user for API key
	fmt.Println("Please enter your BF4DB Patreon API key. You can get one here: https://bf4db.com/patreon")
	reader := bufio.NewReader(os.Stdin)
	apiKey, _ = reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey) // Remove whitespace characters from input

	// check if API Key is a 64-character string
	if len(apiKey) != 64 {
		fmt.Println("API key is invalid. Please try again.")
		setApiKey()
	}
	// Save the API as env. variable:
	toWirte := map[string]string{
		"BF4DB_API_KEY": apiKey,
	}
	err := godotenv.Write(toWirte, envPath)
	if err != nil {
		fmt.Println("Error saving API key to .env file", err)
	}
	os.Exit(0)
}

func GlobalSearch(player string) {
	var myUrl string
	if IsDiscord {
		myUrl = fmt.Sprint("https://bf4db.com/api/player/", player, "discordAccount/discord?api_token=", apiKey)
	} else {
		myUrl = fmt.Sprint("https://bf4db.com/api/player/", player, "/search?api_token=", apiKey)
	}
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, myUrl, nil)
	if err != nil {
		fmt.Println("NewRequest Error")
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Request Error (client.Do):", err)
		return
	}

	for retryCount := 0; retryCount < 3; {
		if res.StatusCode == http.StatusTooManyRequests {
			fmt.Println("Too Many Requests, waiting 1 minute | Status Code:", res.StatusCode, res.Status)
			retryCount++
			time.Sleep(60 * time.Second)
			retryCount++
			res, err = client.Do(req)
		} else {
			break
		}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("ReadAll Error")
		return
	}

	// Use PlayerData interface to unmarshal the JSON response from BF4DB
	var bfdbApi PlayerData
	if IsDiscord {
		bfdbApi = &BFDiscord{} // Discord sends a [PlayerId] instead of [ID]
	} else {
		bfdbApi = &BFDBAPI{}
	}

	err = json.Unmarshal(body, &bfdbApi)
	if err != nil {
		fmt.Println("Your API key is invalid! Please set a new one with bf4db -c")
		if IsDebug { // if debug is enabled, print the response body
			fmt.Println(string(body))
		}
		return
	}
	if len(bfdbApi.GetData()) == 0 {
		if IsDebug {
			fmt.Println(bfdbApi.GetData(), "No player found") // For debug only
		}
		return
	}
	// print number of players founds when > 15 (as it is harder to read)
	if len(bfdbApi.GetData()) > 15 {
		fmt.Println("More than 15 players found! Total of", len(bfdbApi.GetData()), "\n")
	}
	for x := range bfdbApi.GetData() {
		if bfdbApi.GetData()[x].BanReason == "" {
			bfdbApi.GetData()[x].BanReason = "Under review"
		}

		if IsDebug == true { // For debug only
			fmt.Println("Received Data:", bfdbApi.GetData()[x])
			continue
		}
		// if is nil, do nothing
		id := bfdbApi.GetData()[x].ID // if Global Search ID is 0, check the PlayerId (Discord)
		if id == 0 {
			id = bfdbApi.GetData()[x].PlayerId
		}
		if id == 0 {
			continue
		}

		bfdbURL := fmt.Sprint("https://bf4db.com/player/", id, "/")
		bf4crURL := fmt.Sprint("http://bf4cheatreport.com/?pid=", id, "&uid=&cnt=200&startdate=", time.Now().Format("200601021504"))
		bfAgency := fmt.Sprint("https://battlefield.agency/player/by-persona_id/bf4/=", id)
		fmt.Printf("%v | %v | Cheat score = %v | %v\n Cheat Report: %v\n BF Agency: %v\n\n",
			bfdbApi.GetData()[x].Name, bfdbApi.GetData()[x].BanReason, bfdbApi.GetData()[x].CheatScore, bfdbURL, bf4crURL, bfAgency)
	}
}

type PlayerData interface {
	GetData() []Player
}

type Player struct {
	ID         int       `json:"id"`
	PlayerId   int       `json:"player_id"`
	Name       string    `json:"name"`
	IsBanned   int       `json:"is_banned"`
	BanReason  string    `json:"ban_reason"`
	EaGuid     string    `json:"ea_guid"`
	PbGuid     string    `json:"pb_guid"`
	CheatScore int       `json:"cheat_score"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BFDBAPI struct {
	Data []Player `json:"data"`
}

func (b BFDBAPI) GetData() []Player {
	return b.Data
}

type BFDiscord struct {
	Data []Player `json:"data"`
}

func (b BFDiscord) GetData() []Player {
	return b.Data
}
