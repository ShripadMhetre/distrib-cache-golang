package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shripadmhetre/distrib-cache-golang/cache"
)

func main() {
	cache := cache.New()

	// SET key1 val1 EX 15
	cache.Set([]byte("key1"), []byte("val1"), time.Duration(15*time.Second))

	// GET key1
	val, err := cache.Get([]byte("key1"))

	if err != nil {
		log.Fatal("Error: ", err)
	}

	fmt.Printf("key: key1, value: %s", string(val))

	// EXISTS key1
	isExist := cache.Exists([]byte("key1"))

	fmt.Println("\nDoes cache has key: key1 =>", isExist)

	// Testing the key expiration logic
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		<-time.After(15 * time.Second)
		val, err := cache.Get([]byte("key1"))

		if err != nil {
			log.Fatal("Goroutine Error: ", err)
		}

		fmt.Printf("key: key1, value: %s", string(val))

		wg.Done()
	}()

	wg.Wait()
}
