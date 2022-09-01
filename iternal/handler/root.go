package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"l0/iternal/service"
)

func NewRoot(so *service.Order) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUid := r.URL.Query().Get("order_uid")
		b, ok := so.GetByUid(orderUid)
		fmt.Println("new root response:", string(b))
		if !ok {
			//w.WriteHeader(http.StatusNotFound)
			//return
		}
		w.WriteHeader(http.StatusOK)
		http.ServeFile(w, r, "./web/index.html")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(b)
		if err != nil {
			return
		}
	}
}
