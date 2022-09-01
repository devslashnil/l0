package main

import (
	"fmt"
	"l0/iternal/handler"
	"net/http"
	"os"
	"time"

	"l0/iternal/repository"
	"l0/iternal/service"
	"l0/iternal/sub"
	"l0/iternal/util"

	"github.com/patrickmn/go-cache"
)

func main() {
	util.LoadEnv(".env")
	r := repository.NewOrderRepoFromUrl(os.Getenv("DATABASE_URL"))
	c := cache.New(5*time.Minute, 10*time.Minute)
	util.InitCache(c, r)
	so := service.NewOrderService(c, r)
	sc := sub.NewStanConn()
	sub.Subscribe(sc, "order", sub.NewOrderMsgHandler(so))
	http.HandleFunc("/", handler.NewRoot(so))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./web/assets"))))
	fmt.Printf("Starting server at port 8080\n")
	http.ListenAndServe(":8080", nil)
}
