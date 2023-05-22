package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func main() {
	db, err := sql.Open("mysql", "salman:salman_pass@tcp(localhost)/perpustakaan")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO pengguna_ordinary(nama, golongan_darah, tinggi_badan, usia) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(fmt.Sprintf("%s %d", uuid.New(), i), random_golongan_darah(), random_tinggi_badan(), random_usia())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func random_usia() int {
	rand.Seed(time.Now().UnixNano())

	// Define the range
	min := 22
	max := 56

	// Generate a random number within the range
	return rand.Intn(max-min+1) + min
}

func random_tinggi_badan() int {
	rand.Seed(time.Now().UnixNano())

	// Define the range
	min := 155
	max := 185

	// Generate a random number within the range
	return rand.Intn(max-min+1) + min
}

func random_golongan_darah() int {
	rand.Seed(time.Now().UnixNano())
	var (
		golongan_darah = []int{
			0, // o
			1, // a
			2, // b
			3, // ab
		}
	)

	randomIndex := rand.Intn(len(golongan_darah))

	return golongan_darah[randomIndex]
}
