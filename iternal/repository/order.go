package repository

import (
	"context"
	"fmt"
	"os"

	"l0/iternal/model"

	"github.com/jackc/pgx/v4"
)

type Order struct {
	conn *pgx.Conn
}

func NewOrderRepoFromUrl(url string) *Order {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &Order{conn}
}

func (r *Order) Create(o *model.Order) error {
	fmt.Printf("Creating stuff: %v\n", o)
	query := "CALL add_order($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);"
	err := r.conn.QueryRow(
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
		//return err
	}
	err = r.createDelivery(&o.Delivery, o.OrderUid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		//return err
	}
	err = r.createPayment(&o.Payment, o.OrderUid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		//return err
	}
	for _, item := range o.Items {
		err = r.createItem(&item, o.OrderUid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			//return err
		}
	}
	return nil
}

func (r *Order) createDelivery(d *model.Delivery, OrderUid string) error {
	query := "CALL add_delivery($1, $2, $3, $4, $5, $6, $7, $8);"
	err := r.conn.QueryRow(
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
		fmt.Fprintf(os.Stderr, "createDelivery QueryRow failed: %v\n", err)
		return err
	}
	return nil
}

func (r *Order) createItem(i *model.Item, OrderUid string) error {
	query := "CALL add_order_item($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);"
	err := r.conn.QueryRow(
		context.Background(),
		query,
		OrderUid,
		i.Sale,
		i.ChrtId,
		i.TrackNumber,
		i.Price,
		i.Name,
		i.Rid,
		i.Size,
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

func (r *Order) createPayment(p *model.Payment, OrderUid string) error {
	query := "CALL add_payment($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	err := r.conn.QueryRow(
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

func (r *Order) GetOrder(uid string) (*model.Order, error) {
	var order model.Order
	query := "CALL get_order( $1 );"
	err := r.conn.QueryRow(context.Background(), query, uid).Scan(&order)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return &order, err
	}
	return &order, nil
}

func (r *Order) GetAllOrders() ([]*model.Order, error) {
	query := `SELECT "get_all_orders"();`
	rows, err := r.conn.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get_all_orders failed: %v\n", err)
	}
	orders := make([]*model.Order, 0)
	for rows.Next() {
		var order model.Order
		if err = rows.Scan(&order); err != nil {
			fmt.Fprintf(os.Stderr, "Query Scan failed: %v\n", err)
			return orders, err
		}
		orders = append(orders, &order)
	}
	if err = rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error from iterating over rows: %v\n", err)
		return orders, err
	}
	return orders, nil
}
