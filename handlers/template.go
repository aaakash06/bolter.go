package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["templateId"])
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("something\n cool"))
}
