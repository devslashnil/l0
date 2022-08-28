package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"

	stan "github.com/nats-io/stan.go"
)

const (
	msgNum = 10
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
	err = json.Unmarshal(file, &order)
	if err != nil {
		log.Fatalln(err)
	}
	// Tick return
	orders := populateOrders(order, msgNum)
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	i := 0
	for {
		select {
		case <-ticker.C:
			fmt.Println("foo publish")
			err = sc.Publish("foo", file)
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

func populateOrders(template Order, l int) []Order {
	orders := make([]Order, l)
	for i := 0; i < l; i++ {
		order := mutate(&template, i)
		payment := template.Payment
		order.Payment = mutate(&payment, i)
		items := make([]Item, len(template.Items))
		for _, v := range items {
			item := mutate(&v, i)
			items = append(items, item)
		}
		orders = append(orders, order)
	}
	return orders
}

// mutate change a struct values
func mutate[T any](a *T, num int) T {
	v := reflect.Indirect(reflect.ValueOf(a))
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch prop := field.Interface().(type) {
		case int:
			field.SetInt(int64(prop + rand.Intn(100)))
		case string:
			field.SetString(fmt.Sprintf("test-data-%d:%v", num, prop))
		}
	}
	return *a
}
