package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"l0/iternal/model"
	"l0/iternal/repository"

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
		// При получение невалидных данных обрабатывем их другим хэндлером и не производим сохранение в базу и кэш
		so.handleUnknownMsg(m)
		return
	}
	err = so.r.Create(&o)
	if err != nil {
		log.Fatal(err)
		// Чтобы не потерять данные из-за сбоя базы, продолжаем поток выполнения и сохраняем in-memory
	}
	so.c.Set(o.OrderUid, m.Data, cache.DefaultExpiration)
	fmt.Fprintf(os.Stdout, "Next order_uid saved to util: %v", o.OrderUid)
}

func (so *Order) GetByUid(uid string) (*model.Order, bool) {
	i, ok := so.c.Get(uid)
	if !ok {
		return nil, false
	}
	o, ok := i.(*model.Order)
	if !ok {
		return nil, false
	}
	return o, true
}

func (so *Order) handleUnknownMsg(m *stan.Msg) {
	log.Printf("Got unsuported data from order-pub: %v", m.Data)
}
