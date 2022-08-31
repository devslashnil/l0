package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"l0/iternal/handler"
	"l0/iternal/repository"
	"l0/iternal/service"
	"l0/iternal/sub"
	"l0/iternal/util"

	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
)

func main() {
	godotenv.Load()
	r := repository.NewOrderRepoFromUrl(os.Getenv("DATABASE_URL"))
	c := cache.New(5*time.Minute, 10*time.Minute)
	util.InitCache(c, r)
	so := service.NewOrderService(c, r)
	sc := sub.NewStanConn()
	(*sc).Subscribe("order", sub.NewOrderMsgHandler(so))
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.NewRootHandler(so))
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
