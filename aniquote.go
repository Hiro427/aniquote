package main

// {
// 	status: `success`,
// 	data: {
// 		content: `Actually... Ponytails turn me on... that ponytail you had back then, it looked so good it was criminal!`,
// 		anime: {
// 			id: 319,
// 			name: `The Melancholy of Haruhi Suzumiya`
// 		},
// 		character: {
// 			id: 401,
// 			name: `Kyon`
// 		}
// 	}
// }
import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mitchellh/go-wordwrap"

	_ "github.com/mattn/go-sqlite3"
)

type Quote struct {
	Content   string
	Anime     string
	Character string
}

func api(db *sql.DB) {
	url := "https://animechan.io/api/v1/quotes/random"
	resp, _ := http.Get(url)
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	data := result["data"].(map[string]interface{})

	quote := Quote{
		Content:   data["content"].(string),
		Anime:     data["anime"].(map[string]interface{})["name"].(string),
		Character: data["character"].(map[string]interface{})["name"].(string),
	}

	insertDB(db, quote.Content, quote.Anime, quote.Character)
}
func makeDB() *sql.DB {
	homeDir, _ := os.UserHomeDir()
	dbPath := homeDir + "/.quotes.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS quotes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT,
        anime TEXT,
        character TEXT
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func quoteExists(db *sql.DB, content string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM quotes WHERE content = ? LIMIT 1)`
	err := db.QueryRow(query, content).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

// Function to insert a quote into the database, only if it's not a duplicate
func insertDB(db *sql.DB, quote string, anime string, character string) {
	if !quoteExists(db, quote) {
		insertQuoteSQL := `INSERT INTO quotes (content, anime, character) VALUES (?, ?, ?)`
		_, err := db.Exec(insertQuoteSQL, quote, anime, character)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Quote Inserted")
	} else {
		fmt.Println("Quote already exists, skipping insert.")
	}
}

func insertCustomDB(db *sql.DB) {
	type CustomQuote struct {
		Content     string
		Description string
		Person      string
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the quote: ")
	customQuote, _ := reader.ReadString('\n')
	customQuote = strings.TrimSpace(customQuote)
	fmt.Println("Enter Description (book, poem, speech etc): ")
	customDescription, _ := reader.ReadString('\n')
	customDescription = strings.TrimSpace(customDescription)
	fmt.Println("Enter Name: ")
	customName, _ := reader.ReadString('\n')
	customName = strings.TrimSpace(customName)

	customq := CustomQuote{
		Content:     customQuote,
		Description: customDescription,
		Person:      customName,
	}
	insertDB(db, customq.Content, customq.Description, customq.Person)

}

// Function to retrieve a random quote from the database
func getRandomQuote(db *sql.DB) (string, string, string) {
	var content, anime, character string
	query := `SELECT content, anime, character FROM quotes ORDER BY RANDOM() LIMIT 1`
	err := db.QueryRow(query).Scan(&content, &anime, &character)
	if err != nil {
		log.Fatal(err)
	}
	return content, anime, character
}

func main() {
	db := makeDB()
	defer db.Close()

	if len(os.Args) < 2 {
		fmt.Println("Please provide an option: 'update' or 'random'")
		return
	}

	opt := os.Args[1]

	if opt == "update" {
		for i := 0; i < 10; i++ {
			api(db)
		}
	} else if opt == "random" {
		content, anime, character := getRandomQuote(db)
		wrappedQuote := wordwrap.WrapString(content, 100)
		fmt.Printf("%s\n \n-%s, %s\n", wrappedQuote, character, anime)
	} else if opt == "insert" {
		insertCustomDB(db)
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter another quote? (y/n): ")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)
		if option == "y" {
			insertCustomDB(db)
		}
	} else {
		fmt.Println("Invalid option. Use 'update' to fetch new quotes or 'random' to get a random quote.")
	}
}
