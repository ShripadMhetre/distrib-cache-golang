package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shripadmhetre/distrib-cache-golang/cache"
)

func main() {
	// var cache cache.Cache

	cache := cache.New()

	cache.Set([]byte("key1"), []byte("val1"), time.Duration(15*time.Second))
	val, err := cache.Get([]byte("key1"))

	if err != nil {
		log.Fatal("Error: ", err)
	}

	fmt.Printf("key: key1, value: %s", string(val))
	var wg sync.WaitGroup

	fmt.Println("\nDoes cache has key: key1? ", cache.Has([]byte("key1")))

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
