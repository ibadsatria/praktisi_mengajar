package main

// Assuming you have imported the following packages
import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/patrickmn/go-cache"
)

var c *cache.Cache

// Define a struct for the pengguna table
type Pengguna struct {
	ID            int    `json:"id"`
	Nama          string `json:"nama"`
	GolonganDarah int    `json:"golongan_darah"`
	TinggiBadan   int    `json:"tinggi_badan"`
	Usia          int    `json:"usia"`
}

func main() {
	c = cache.New(5*time.Minute, 10*time.Minute) // Adjust the expiration and cleanup intervals as needed
	// Create a reader to read input from standard input
	reader := bufio.NewReader(os.Stdin)

	// Loop until the user enters "exit"
	for {
		// Prompt the user to enter a command
		fmt.Print("Enter a command: ")

		// Read the input until a newline
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Trim the input and convert it to lower case
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)

		// Check if the user wants to exit
		if input == "exit" {
			fmt.Println("Bye!")
			break
		}

		id, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("don't understand. Bye!")
			break
		}

		pengguna, err := getPengguna(context.Background(), id)
		if err != nil {
			fmt.Errorf("error %v", err)
			return
		}

		printPengguna(pengguna)
	}

}

func printPengguna(p Pengguna) {
	fmt.Printf("Pengguna\nid: %d\nnama: %s\nusia: %d\ntinggi badan: %dcm\ngolongan darah: %s\n=========\n", p.ID, p.Nama, p.Usia, p.TinggiBadan, parsegolongandarah(p.GolonganDarah))
}

func parsegolongandarah(g int) string {
	switch g {
	case 0:
		return "o"
	case 1:
		return "a"
	case 2:
		return "b"
	case 3:
		return "ab"
	}
	return ""
}

// Define a function that takes an id and returns a Pengguna object
func getPengguna(ctx context.Context, id int) (Pengguna, error) {
	// Create a Pengguna variable to store the result
	var p Pengguna

	// Try to get the data from memory cache first

	key := fmt.Sprintf("pengguna:%d", id) // Use a unique key for each id key: pengguna:2304109
	val, found := c.Get(key)
	if found {
		fmt.Println("is fetching from memory cache")
		// Cache hit, unmarshal the value to p
		err := json.Unmarshal(val.([]byte), &p)
		if err != nil {
			return p, err
		}
		return p, nil
	}

	// Cache miss, try to get the data from redis next
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380", // Adjust the address and password as needed
		Password: "",
		DB:       0,
	})
	val, err := rdb.Get(ctx, key).Result()
	if err == nil {
		fmt.Println("is fetching from memory redis (remote cache)")
		// Redis hit, unmarshal the value to p
		err := json.Unmarshal([]byte(fmt.Sprint(val)), &p)
		if err != nil {
			return p, err
		}
		// Set the value to memory cache for future use
		c.Set(key, []byte(fmt.Sprint(val)), cache.DefaultExpiration)
		return p, nil
	}

	// Redis miss, try to get the data from mysql last
	db, err := sql.Open("mysql", "salman:salman_pass@tcp(localhost)/perpustakaan")
	if err != nil {
		return p, err
	}
	defer db.Close()
	row := db.QueryRow("SELECT id, nama, golongan_darah, tinggi_badan, usia FROM pengguna WHERE id = ?", id)
	err = row.Scan(&p.ID, &p.Nama, &p.GolonganDarah, &p.TinggiBadan, &p.Usia)
	if err != nil {
		return p, err
	}
	// Mysql hit, marshal the value to JSON and set it to redis and memory cache for future use
	valBytes, err := json.Marshal(p)
	if err != nil {
		return p, err
	}

	fmt.Println("is fetching from database as last failover")
	rdb.Set(ctx, key, valBytes, time.Hour) // Adjust the expiration time as needed
	c.Set(key, valBytes, 10*time.Second)   // set expiry time to 10 seconds
	return p, nil
}
