package util

import (
	"fmt"
	"l0/iternal/repository"
	"os"

	"github.com/patrickmn/go-cache"
)

func InitCache(c *cache.Cache, r *repository.OrderRepo) *cache.Cache {
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
