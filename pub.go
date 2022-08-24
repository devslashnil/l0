package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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
	var order Order
	err := json.Unmarshal(file, &order)
	// Tick return
	//ticker := time.NewTicker(10 * time.Second)
	//quit := make(chan struct{})
	//for {
	//	select {
	//	case <-ticker.C:
	//		fmt.Println("foo")
	//		err = sc.Publish("foo", file)
	//		if err != nil {
	//			log.Fatalln(err)
	//		}
	//	case <-quit:
	//		ticker.Stop()
	//		return
	//	}
	//}
	fmt.Println("end")
}

func populateOrders(template Order, len int) []Order {
	orders := make([]Order, 0)
	for i := 0; i < len; i++ {
		order := template
		payment := template.Payment
		delivery := template.Delivery
		items := template.Items
	}
	return orders
}
