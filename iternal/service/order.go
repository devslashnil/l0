package service

import (
	"encoding/json"
	"fmt"
	"l0/iternal/model"
	"l0/iternal/repository"
	"log"
	"os"

	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
)

type Order struct {
	r *repository.Order
	c *cache.Cache
}

func NewOrderService(c *cache.Cache, r *repository.Order) *Order {
	return &Order{r, c}
}

func (so *Order) SaveFromMsg(m *stan.Msg) {
	var o model.Order
	err := json.Unmarshal(m.Data, &o)
	if err != nil {
		so.handleUnknownMsg(m)
	}
	err = so.r.Create(&o)
	if err != nil {
		// todo: gracefully handle
		log.Fatal(err)
		return
	}
	so.c.Set(o.OrderUid, m.Data, cache.DefaultExpiration)
	fmt.Fprintf(os.Stdout, "Next order_uid saved to util: %v", o.OrderUid)
}

func (so *Order) GetByUid(uid string) ([]byte, bool) {
	o, ok := so.c.Get(uid)
	if !ok {
		return nil, ok
	}
	b, ok := o.([]byte)
	if !ok {
		log.Fatalf("If there isn't []byte in cache than what is it: %T?", o)
		return nil, ok
	}
	return b, true
}

func (so *Order) handleUnknownMsg(m *stan.Msg) {
	// todo: cry loudly
}
