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
	// load cache
	orders, err := getAllOrders(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get all orders: %v", err)
	}
	for _, order := range orders {
		c.Set(order.OrderUid, order, cache.DefaultExpiration)
	}
	fmt.Fprintf(os.Stdout, "Cache recovered")

	sc.Subscribe("pub-order", func(m *stan.Msg) {
		var order Order
		err = json.Unmarshal(m.Data, &order)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Got invalid data order from pub")
			return
		}
		err = order.Create(conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create order_uid: %v", order.OrderUid)
			return
		}
		c.Set(order.OrderUid, string(m.Data), cache.DefaultExpiration)
		fmt.Fprintf(os.Stdout, "Next order_uid saved to cache: %v", order.OrderUid)
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		orderUid := r.URL.Query().Get("order_uid")
		order, ok := c.Get(orderUid)
		if ok {
			fmt.Fprintf(w, "%v", order)
			return
		}
		order, err = getOrder(conn, orderUid)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)

		c.Set(orderUid, order, cache.DefaultExpiration)
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}

func (o Order) Create(conn *pgx.Conn) error {
	query := "CALL add_order($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);"
	err := conn.QueryRow(
		context.Background(),
		query,
		o.OrderUid,
		o.TrackNumber,
		o.Entry,
		o.Locale,
		o.InternalSignature,
		o.CustomerId,
		o.DeliveryService,
		o.Shardkey,
		o.SmId,
		o.DateCreated,
		o.OofShard,
	).Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return err
	}
	err = o.Delivery.Create(conn, o.OrderUid)
	if err != nil {
		return err
	}
	err = o.Payment.Create(conn, o.OrderUid)
	if err != nil {
		return err
	}
	for _, item := range o.Items {
		err = item.Create(conn, o.OrderUid)
		if err != nil {
			return err
		}
	}
	// TODO: do rollback on error
	return nil
}

func (d Delivery) Create(conn *pgx.Conn, OrderUid string) error {
	query := "CALL add_payment($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);"
	err := conn.QueryRow(
		context.Background(),
		query,
		OrderUid,
		d.Name,
		d.Phone,
		d.Zip,
		d.City,
		d.Address,
		d.Region,
		d.Email,
	).Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return err
	}
	return nil
}

func (i Item) Create(conn *pgx.Conn, OrderUid string) error {
	query := "CALL add_order_item($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);"
	err := conn.QueryRow(
		context.Background(),
		query,
		OrderUid,
		i.Sale,
		i.ChrtId,
		i.TrackNumber,
		i.Price,
		i.Rid,
		i.Name,
		i.TotalPrice,
		i.NmId,
		i.Brand,
		i.Status,
	).Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return err
	}
	return nil
}

func (p Payment) Create(conn *pgx.Conn, OrderUid string) error {
	query := "CALL add_payment($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	err := conn.QueryRow(
		context.Background(),
		query,
		p.Transaction,
		OrderUid,
		p.RequestId,
		p.Currency,
		p.Provider,
		p.Amount,
		p.PaymentDt,
		p.Bank,
		p.DeliveryCost,
		p.GoodsTotal,
		p.CustomFee,
	).Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return err
	}
	return nil
}

func getOrder(conn *pgx.Conn, uid string) (Order, error) {
	var order Order
	query := "CALL get_order( $1 );"
	err := conn.QueryRow(context.Background(), query, uid).Scan(&order)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return order, err
	}
	return order, nil
}

func getAllOrders(conn *pgx.Conn) ([]Order, error) {
	query := "CALL get_all_orders();"
	rows, err := conn.Query(context.Background(), query)
	orders := make([]Order, 0)
	for rows.Next() {
		var order Order
		if err = rows.Scan(&order); err != nil {
			fmt.Fprintf(os.Stderr, "Query Scan failed: %v\n", err)
			return orders, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error from iterating over rows: %v\n", err)
		return orders, err
	}
	return orders, nil
}

type Order struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
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

type Item struct {
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
