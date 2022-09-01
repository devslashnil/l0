package sub

import (
	"fmt"
	"log"
	"os"

	"l0/iternal/service"

	"github.com/nats-io/stan.go"
)

func NewStanConn() *stan.Conn {
	sc, err := stan.Connect("test-cluster", "order-sub")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to nats-streaming-server: %v\n", err)
		os.Exit(1)
	}
	return &sc
}

func NewOrderMsgHandler(so *service.Order) func(*stan.Msg) {
	return func(m *stan.Msg) {
		so.SaveFromMsg(m)
	}
}

func Subscribe(sc *stan.Conn, subject string, msgHandler func(*stan.Msg)) {
	_, err := (*sc).Subscribe(subject, msgHandler)
	if err != nil {
		log.Fatal(err)
	}
}
