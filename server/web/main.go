package main

import (
	"github.com/angelofallars/htmx-go"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", helloWorld)
	http.ListenAndServe(":8080", nil)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	htmx.NewResponse().
		Retarget("#name").
		RenderTempl(r.Context(), w, Page())
}
