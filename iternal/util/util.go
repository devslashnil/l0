package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"l0/iternal/repository"

	"github.com/patrickmn/go-cache"
)

func InitCache(c *cache.Cache, r *repository.Order) *cache.Cache {
	orders, err := r.GetAllOrders()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get all orders: %v", err)
	}
	for _, order := range orders {
		fmt.Println(order)
		c.Set(order.OrderUid, order, cache.DefaultExpiration)
	}
	fmt.Fprintf(os.Stdout, "Cache recovered\n")
	return c
}

func LoadEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		keyAndVal := strings.Split(sc.Text(), "=")
		os.Setenv(keyAndVal[0], keyAndVal[1])
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
}
