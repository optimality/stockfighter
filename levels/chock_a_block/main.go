package main

import (
	"fmt"

	"github.com/optimality/stockfighter"
)

const (
	InstanceID = 28574
)

func main() {
	client, err := stockfighter.NewStockfighterClient()
	if err != nil {
		panic(err)
	}

	instance, err := client.RestartInstance(InstanceID)
	if err != nil {
		panic(err)
	}
	venue := instance.Venues[0]
	stock := instance.Tickers[0]
	account := instance.Account

	quote, err := client.Quote(venue, stock)
	if err != nil {
		panic(err)
	}
	for quote.Last == 0 {
		quote, err = client.Quote(venue, stock)
		if err != nil {
			panic(err)
		}
	}
	price := quote.Last
	fmt.Printf("Buying at %v\n", price)
	fmt.Println()
	bought := 0
	for bought < 100000 {
		fmt.Printf("\rBought %v / 100000", bought)
		status, err := client.PostOrder(stockfighter.PostOrderRequest{
			Account:   account,
			Venue:     venue,
			Stock:     stock,
			Price:     price,
			Qty:       100000 - bought,
			Direction: stockfighter.Buy,
			OrderType: stockfighter.ImmediateOrCancel,
		})
		if err != nil {
			panic(err)
		}
		bought += status.TotalFilled
	}
}
