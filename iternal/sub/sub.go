package sub

import (
	"fmt"
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
		// todo: log instead print
		fmt.Printf("Got message: %s", m.Data)
		so.SaveFromMsg(m)
	}
}
