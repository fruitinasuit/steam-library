package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// Steam API response structures
type PlayerSummaryResponse struct {
	Response struct {
		Players []struct {
			Personaname string `json:"personaname"`
		} `json:"players"`
	} `json:"response"`
}

type OwnedGamesResponse struct {
	Response struct {
		Games []map[string]interface{} `json:"games"`
	} `json:"response"`
}

// getSteamPersonaname fetches the user's display name from Steam API
func getSteamPersonaname(steamID, apiKey string) (string, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s&format=json", apiKey, steamID)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data PlayerSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if len(data.Response.Players) == 0 {
		return "", fmt.Errorf("no player found")
	}

	return data.Response.Players[0].Personaname, nil
}

// getSteamLibrary fetches the user's owned games from Steam API
func getSteamLibrary(steamID, apiKey string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=1", apiKey, steamID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data OwnedGamesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response.Games, nil
}

// outputGamesToCSV writes the games data to a CSV file
func outputGamesToCSV(games []map[string]interface{}, personaname string) error {
	if len(games) == 0 {
		return fmt.Errorf("no games to output")
	}

	filename := fmt.Sprintf("steam_library_%s.csv", personaname)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Get all unique keys for headers
	keySet := make(map[string]bool)
	for _, game := range games {
		for key := range game {
			keySet[key] = true
		}
	}

	// Create headers
	var headers []string
	for key := range keySet {
		headers = append(headers, key)
	}

	// Write headers
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write data rows
	for _, game := range games {
		var row []string
		for _, header := range headers {
			if val, ok := game[header]; ok {
				row = append(row, formatValue(val))
			} else {
				row = append(row, "")
			}
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// formatValue converts interface{} to string for CSV
func formatValue(val interface{}) string {
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		// For complex types, convert to string representation
		return fmt.Sprintf("%v", v)
	}
}

func main() {
	// Load credentials from environment variables
	apiKey := os.Getenv("STEAM_API_KEY")
	steamID := os.Getenv("STEAM_ID")

	if apiKey == "" {
		fmt.Println("Error: STEAM_API_KEY environment variable is not set")
		return
	}

	if steamID == "" {
		fmt.Println("Error: STEAM_ID environment variable is not set")
		return
	}

	// Get personaname
	personaname, err := getSteamPersonaname(steamID, apiKey)
	if err != nil {
		fmt.Printf("Error getting personaname: %v\n", err)
		return
	}

	// Get steam library
	games, err := getSteamLibrary(steamID, apiKey)
	if err != nil {
		fmt.Printf("Error getting steam library: %v\n", err)
		return
	}

	// Output to CSV
	if err := outputGamesToCSV(games, personaname); err != nil {
		fmt.Printf("Error outputting to CSV: %v\n", err)
		return
	}

	fmt.Printf("Successfully downloaded Steam library for %s\n", personaname)
}
