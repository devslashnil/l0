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
		order := template

		payment := template.Payment
		delivery := template.Delivery
		items := template.Items
	}
	return orders
}

// Mutate mutate a struct, where a - pointer
func Mutate(a any) {
	v := reflect.Indirect(reflect.ValueOf(a))
	fmt.Printf("%v\n", v)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch prop := field.Interface().(type) {
		case int:
			field.SetInt(int64(prop + rand.Intn(100)))
		case string:
			// todo: replace i on j
			field.SetString(fmt.Sprintf("test-data-%d:%v", i, field.String()))
		case []Item:
			fmt.Println("array")
		default:
			//fmt.Printf("prop %v\n", prop)
			l := reflect.Indirect(reflect.ValueOf(&prop)).Elem() //.NumField()
			fmt.Printf("l %v\n", l)
			//Mutate(&prop)
		}
	}
}
