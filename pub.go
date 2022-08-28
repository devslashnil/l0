package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"

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
		order := mutate(template)

		payment := template.Payment
		delivery := template.Delivery
		items := template.Items
	}
	return orders
}

// mutate change a struct values
func mutate(a any, num int) any {
	v := reflect.ValueOf(a)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch prop := field.Interface().(type) {
		case int:
			field.SetInt(int64(prop + rand.Intn(100)))
		case string:
			field.SetString(fmt.Sprintf("test-data-%d:%v", num, field.String()))
		}
	}
	return a
}
