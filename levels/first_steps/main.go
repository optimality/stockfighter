package main

import (
	"log"

	"github.com/optimality/stockfighter"
)

func main() {
	venue := "GFGBEX"
	stock := "RYYA"
	account := "LPS79640982"

	client, err := stockfighter.NewStockfighterClient()
	if err != nil {
		panic(err)
	}
	status, err := client.PostOrder(stockfighter.PostOrderRequest{
		Account:   account,
		Venue:     venue,
		Stock:     stock,
		Price:     0,
		Qty:       100,
		Direction: stockfighter.Buy,
		OrderType: stockfighter.Market,
	})
	if err != nil {
		panic(err)
	}
	for status.Open {
		status, err = client.Order(venue, stock, status.ID)
		if err != nil {
			panic(err)
		}
	}
	log.Printf("Status: %+v\n", status)
}
