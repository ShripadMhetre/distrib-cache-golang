package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/shripadmhetre/distrib-cache-golang/cache"
	"github.com/shripadmhetre/distrib-cache-golang/client"
)

func main() {
	listenAddr := flag.String("listenaddr", "localhost:3000", "listen address of server")
	leaderAddr := flag.String("leaderaddr", "", "listen address of the leader")

	flag.Parse()

	options := ServerOptions{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(time.Second * 10)
		if options.IsLeader {
			SimulateClient()
		}
	}()

	server := NewServer(options, cache.New())
	server.Run()
}

// client simulation function
func SimulateClient() {
	client, err := client.New("localhost:3000")
	for i := 0; i < 10; i++ {
		go func(i int) {
			if err != nil {
				log.Fatal("Error connecting to server: ", err)
			}

			var (
				key   = []byte(fmt.Sprintf("key_%d", i))
				value = []byte(fmt.Sprintf("val_%d", i))
			)

			// SET key to value
			err = client.Set(context.Background(), key, value, 0)
			if err != nil {
				log.Fatal(err)
			}

			// GET key
			fetchedValue, err := client.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Value: ", string(fetchedValue))

			client.Close()
		}(i)
	}

	// Exists key
	isExists, err := client.Exists(context.Background(), []byte("key_2"))
	fmt.Printf("Is Key: %s exist? %s", []byte("key_2"), isExists)

}
