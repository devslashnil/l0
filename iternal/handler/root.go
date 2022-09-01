package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"l0/iternal/service"
)

func NewRoot(so *service.Order) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUid := r.URL.Query().Get("order_uid")
		o, ok := so.GetByUid(orderUid)
		tmpl, err := template.ParseFiles("./web/index.html")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		if !ok {
			tmpl.Execute(w, struct {
				Data string
			}{"По такому id ничего не найдено"})
			return
		}
		b, _ := json.MarshalIndent(*o, "", "\t")
		tmpl.Execute(w, struct {
			Data string
		}{Data: string(b)})
	}
}
