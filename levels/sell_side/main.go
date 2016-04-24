package main

import (
	"fmt"
	"time"

	"github.com/optimality/stockfighter"
)

const (
	InstanceID = 28597
	Alpha      = .0
)

type Trader struct {
	Client        stockfighter.StockfighterClient
	Venue         string
	Stock         string
	Account       string
	Position      int
	Bid           stockfighter.OrderStatus
	BidTimestamp  time.Time
	Ask           stockfighter.OrderStatus
	AskTimestamp  time.Time
	Cash          int
	PriceEstimate int
}

func (t *Trader) Run() {
	var oldQuote stockfighter.QuoteResponse
	for {
		t.UpdateOrderStatus()

		quote, err := t.Client.Quote(t.Venue, t.Stock)
		if err != nil {
			panic(err)
		}

		nav := t.Position*quote.Last + t.Cash
		fmt.Printf("Cash: $%v Position: %v NAV: $%v "+
			"OurBid: %v @ $%v  OurAsk: %v @ $%v Last: %v\n",
			float64(t.Cash)/100, t.Position, float64(nav)/100,
			t.Bid.OriginalQty-t.Bid.TotalFilled, float64(t.Bid.Price)/100,
			t.Ask.OriginalQty-t.Ask.TotalFilled, float64(t.Ask.Price)/100, float64(quote.Last)/100)

		if oldQuote.Last != quote.Last {
			t.UpdatePriceEstimate(quote.Last)
		}
		oldQuote = quote
		t.UpdateBid()
		t.UpdateAsk()
	}
}

func (t *Trader) UpdatePriceEstimate(lastPrice int) {
	if t.PriceEstimate == 0 {
		t.PriceEstimate = lastPrice
	} else {
		t.PriceEstimate = int(Alpha*float64(t.PriceEstimate) + (1.0-Alpha)*float64(lastPrice))
	}
}

func (t *Trader) PositionWeight() float64 {
	return float64(t.Position) / 900.0
}

func (t *Trader) UpdateBid() {
	price := int(float64(t.PriceEstimate) * (.94 - t.PositionWeight()*.04))
	// If our current bid is out of date, kill it.
	if t.Bid.Price != price && t.Bid.ID != 0 {
		bid, err := t.Client.Cancel(t.Venue, t.Stock, t.Bid.ID)
		if err != nil {
			panic(err)
		}
		t.Bid = bid.OrderStatus
		t.UpdateOrderStatus()
	}

	qty := (900 - t.Position) / 4
	if !t.Bid.Open && qty > 0 && price != 0 {
		bid, err := t.Client.PostOrder(stockfighter.PostOrderRequest{
			Account:   t.Account,
			Venue:     t.Venue,
			Stock:     t.Stock,
			Price:     price,
			Qty:       qty,
			Direction: stockfighter.Buy,
			OrderType: stockfighter.Limit,
		})
		if err != nil {
			panic(err)
		}
		t.Bid = bid.OrderStatus
		t.BidTimestamp = time.Time{}
		t.UpdateOrderStatus()
	}
}

func (t *Trader) UpdateAsk() {
	// If our current bid is out of date, kill it.
	price := int(float64(t.PriceEstimate) * (1.04 - t.PositionWeight()*.04))
	if t.Ask.Price != price && t.Ask.ID != 0 {
		ask, err := t.Client.Cancel(t.Venue, t.Stock, t.Ask.ID)
		if err != nil {
			panic(err)
		}
		t.Ask = ask.OrderStatus
		t.UpdateOrderStatus()
	}

	qty := (t.Position + 900) / 4
	if !t.Ask.Open && qty > 0 && price != 0 {
		ask, err := t.Client.PostOrder(stockfighter.PostOrderRequest{
			Account:   t.Account,
			Venue:     t.Venue,
			Stock:     t.Stock,
			Price:     price,
			Qty:       qty,
			Direction: stockfighter.Sell,
			OrderType: stockfighter.Limit,
		})
		if err != nil {
			panic(err)
		}
		t.Ask = ask.OrderStatus
		t.AskTimestamp = time.Time{}
		t.UpdateOrderStatus()
	}
}

func (t *Trader) UpdateOrderStatus() {
	if t.Bid.ID != 0 {
		bidStatus, err := t.Client.Order(t.Venue, t.Stock, t.Bid.ID)
		if err != nil {
			panic(err)
		}
		t.Bid = bidStatus.OrderStatus
		maxTimestamp := time.Time{}
		for _, fill := range bidStatus.Fills {
			if fill.Timestamp.After(t.BidTimestamp) {
				t.Position += fill.Qty
				t.Cash -= fill.Qty * fill.Price
			}
			if fill.Timestamp.After(maxTimestamp) {
				maxTimestamp = fill.Timestamp
			}
		}
		t.BidTimestamp = maxTimestamp
	}

	if t.Ask.ID != 0 {
		askStatus, err := t.Client.Order(t.Venue, t.Stock, t.Ask.ID)
		if err != nil {
			panic(err)
		}
		t.Ask = askStatus.OrderStatus
		maxTimestamp := time.Time{}
		for _, fill := range askStatus.Fills {
			if fill.Timestamp.After(t.AskTimestamp) {
				t.Position -= fill.Qty
				t.Cash += fill.Qty * fill.Price
			}
			if fill.Timestamp.After(maxTimestamp) {
				maxTimestamp = fill.Timestamp
			}
		}
		t.AskTimestamp = maxTimestamp
	}
}

func main() {
	client, err := stockfighter.NewStockfighterClient()
	if err != nil {
		panic(err)
	}

	instance, err := client.RestartInstance(InstanceID)
	if err != nil {
		panic(err)
	}

	trader := Trader{
		Client:  client,
		Venue:   instance.Venues[0],
		Stock:   instance.Tickers[0],
		Account: instance.Account,
	}
	trader.Run()
}
