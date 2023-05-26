package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// db, err := sql.Open("mysql", "salman:salman_pass@tcp(localhost)/perpustakaan")
	db, err := pgxpool.Connect(context.Background(), "postgres://salman:postgres@localhost:5432/pengguna")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// stmt, err := db.Exec("INSERT INTO pengguna(nama, golongan_darah, tinggi_badan, usia) VALUES($1, $2, $3, $4)")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()

	var wg sync.WaitGroup
	wg.Add(10)

	insert_data(db)
	insert_data(db)
	insert_data(db)
	insert_data(db)
	insert_data(db)
	insert_data(db)
	insert_data(db)
	insert_data(db)
	insert_data(db)
	insert_data(db)

	wg.Wait()

}

func insert_data(db *pgxpool.Pool) {
	for i := 0; i < 500000; i++ {
		fmt.Println("inserting data of index ", i)
		_, err := db.Exec(context.Background(),
			"INSERT INTO pengguna(nama, golongan_darah, tinggi_badan, usia) VALUES($1, $2, $3, $4)",
			fmt.Sprintf("%s %d", uuid.New(), i), random_golongan_darah(), random_tinggi_badan(), random_usia())
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
