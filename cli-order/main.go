package main

import (
	"fmt"
	stan "github.com/nats-io/stan.go"
	"time"
)

func main() {
	sc, _ := stan.Connect("test-cluster", "cli-order")

	sub, _ := sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			fmt.Println("waiting on foo")
		case <-quit:
			ticker.Stop()
			return
		}
	}

	sub.Unsubscribe()
}
