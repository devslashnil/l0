package main

import (
	"fmt"
	"log"
	"os"
	"time"

	stan "github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("test-cluster", "pub-order")
	if err != nil {
		log.Fatalln(err)
	}
	file, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatalln(err)
	}
	//data := Order{}
	//_ = json.Unmarshal([]byte(file), &data)
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			fmt.Println("foo")
			err := sc.Publish("foo", file)
			if err != nil {
				log.Fatalln(err)
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
	fmt.Println("end")
}
