package handler

import (
	"encoding/json"
	"l0/iternal/service"
	"net/http"
)

func NewRootHandler(so *service.Order) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUid := r.URL.Query().Get("order_uid")
		b, ok := so.GetByUid(orderUid)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		http.ServeFile(w, r, "index.html")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(b)
		if err != nil {
			return
		}
	}
}
