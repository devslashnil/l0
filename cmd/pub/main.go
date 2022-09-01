package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"

	"l0/iternal/model"

	stan "github.com/nats-io/stan.go"
)

const (
	msgNum = 100
	msgLag = 5 // Sec
)

func main() {
	fmt.Println("Publisher init")
	sc, err := stan.Connect("test-cluster", "order-pub")
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.ReadFile("./api/model.json")
	if err != nil {
		log.Fatal(err)
	}
	var order model.Order
	err = json.Unmarshal(file, &order)
	//fmt.Printf("Unmarshal results: %s\n%v", order, order)
	if err != nil {
		log.Fatal(err)
	}
	// Tick return
	orders := populateOrders(order, msgNum)
	ticker := time.NewTicker(msgLag * time.Second)
	quit := make(chan struct{})
	i := 0
	fmt.Println("1st tick in progress")
	for {
		select {
		case <-ticker.C:
			fmt.Printf("foo publish: %d\n", i)
			b, err := json.Marshal(orders[i])
			//fmt.Printf("json.Marshal(orders[i]): %s\n%v", b, orders[i])
			if err != nil {
				log.Fatalln(err)
			}
			err = sc.Publish("order", b)
			if err != nil {
				log.Fatalln(err)
			}
			i++
			if i >= msgNum {
				quit <- struct{}{}
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

// populateOrders create array of unique model.Order
func populateOrders(template model.Order, l int) []model.Order {
	orders := make([]model.Order, 0, l)
	for i := 0; i < l; i++ {
		order := template
		order = mutate(&order, i)
		order.OrderUid = fmt.Sprintf("%d-i-%v", i, template.OrderUid)
		payment := template.Payment
		order.Payment = mutate(&payment, i)
		items := make([]model.Item, len(template.Items))
		for _, v := range items {
			item := mutate(&v, i)
			items = append(items, item)
		}
		orders = append(orders, order)
	}
	return orders
}

// mutate change a struct of a and add num indexes
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
