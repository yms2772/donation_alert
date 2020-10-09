package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v7/pkg/app"
)

const Address = ":12346"

func main() {
	h := &app.Handler{
		Title: "Donation",
		Styles: []string{
			"/web/app.css",
		},
	}

	log.Println("Running on " + Address)
	if err := http.ListenAndServe(Address, h); err != nil {
		fmt.Println(err)
	}
}
