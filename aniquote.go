package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Quote struct for the API response
type Quote struct {
	Quote     string `json:"quote"`
	Character string `json:"character"`
	Anime     string `json:"anime"`
}

// Save the quote to a new JSON file
func saveQuoteToFile(quoteData Quote, folder string) {
	// Create folder if it doesn't exist
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create a unique filename using the current timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(folder, fmt.Sprintf("quote_%s.json", timestamp))

	// Save the quote to a new JSON file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(quoteData, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	file.Write(jsonData)
}

// Download a random quote from the API and save it
func downloadQuote(folder string) {
	resp, err := http.Get("https://animechan.io/api/v1/quotes/random")
	if err != nil {
		log.Printf("Error fetching quote: %v\n", err)
		time.Sleep(2 * time.Second)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("API is busy or unavailable")
		time.Sleep(2 * time.Second)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		time.Sleep(2 * time.Second)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		return
	}

	data := result["data"].(map[string]interface{})
	quote := data["content"].(string)
	character := data["character"].(map[string]interface{})["name"].(string)
	anime := data["anime"].(map[string]interface{})["name"].(string)

	jsonQuote := Quote{
		Quote:     quote,
		Character: character,
		Anime:     anime,
	}

	// Save the quote
	saveQuoteToFile(jsonQuote, folder)
}

// Select and print a random quote from the folder
func selectRandomQuote(folder string) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatalf("Failed to read directory: %v\n", err)
	}

	var jsonFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			jsonFiles = append(jsonFiles, filepath.Join(folder, file.Name()))
		}
	}

	if len(jsonFiles) == 0 {
		fmt.Println("No quotes available.")
		return
	}

	rand.Seed(time.Now().UnixNano())
	randomFile := jsonFiles[rand.Intn(len(jsonFiles))]

	data, err := ioutil.ReadFile(randomFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v\n", err)
	}

	var quote Quote
	err = json.Unmarshal(data, &quote)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v\n", err)
	}

	fmt.Printf("\"%s\"\n", quote.Quote)
	fmt.Printf("\t%s - \"%s\"\n", quote.Character, quote.Anime)
}

func main() {
	// Define flags
	download := flag.Bool("download", false, "Download a new random quote")
	random := flag.Bool("random", false, "Display a random quote from the saved quotes")

	// Parse the flags
	flag.Parse()

	// Define the folder where quotes are saved
	folder := filepath.Join(os.Getenv("HOME"), ".dotfiles/assets/quotes")

	if *download {
		fmt.Println("Downloading 15 new quotes...")
		for i := 0; i < 15; i++ {
			downloadQuote(folder)
			fmt.Printf("Downloaded quote %d/15\n", i+1)
		}
	}
	if *random {
		selectRandomQuote(folder)
	}

	if !*download && !*random {
		fmt.Println("Please provide a valid option: --download or --random")
	}
}
