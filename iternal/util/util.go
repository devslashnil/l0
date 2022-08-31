package util

import (
	"fmt"
	"os"

	"l0/iternal/repository"

	"github.com/patrickmn/go-cache"
)

func InitCache(c *cache.Cache, r *repository.Order) *cache.Cache {
	orders, err := r.GetAllOrders()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get all orders: %v", err)
	}
	for _, order := range orders {
		c.Set(order.OrderUid, order, cache.DefaultExpiration)
	}
	fmt.Fprintf(os.Stdout, "Cache recovered\n")
	return c
}
