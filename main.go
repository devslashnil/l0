package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
)

// hardcode is bad, but in the task it's convenient for all...
const (
	host     = "localhost"
	port     = 5432
	username = "postgres"
	password = "pass"
	dbname   = "postgres"
)

// TODO: find a way to return whole data at one procedure
// TODO: find a way to insert whole json at one procedure

func main() {
	urlDb := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", username,
		password, host, port, dbname)
	conn, err := pgx.Connect(context.Background(), urlDb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	sc, err := stan.Connect("test-cluster", "sub-order")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to nats-streaming-server: %v\n", err)
		os.Exit(1)
	}
	defer sc.Close()

	c := cache.New(5*time.Minute, 10*time.Minute)

	sc.Subscribe("pub-order", func(m *stan.Msg) {
		var order Order
		err = json.Unmarshal(m.Data, &order)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Got invalid data order from pub")
			return
		}
		err = order.Create()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create order_uid: %v", order.OrderUid)
			return
		}
		c.Set(order.OrderUid, string(m.Data), cache.DefaultExpiration)
		fmt.Fprintf(os.Stdout, "Next order_uid saved to cache: %v", order.OrderUid)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		orderUid := r.URL.Query().Get("order_uid")
		order, err := getOrderByUid(orderUid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to retrive order by uid: %v", orderUid)
			return
		}
		fmt.Fprintf(w, "%v", order)
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

func getOrderByUid(uid string) (Order, error) {
	return Order{}, nil
}

func (o Order) Create() error {
	return nil
}

type Order struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Items   `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Items struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
