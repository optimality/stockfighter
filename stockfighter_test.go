package stockfighter

import "testing"

const (
	TestVenue   = "TESTEX"
	TestStock   = "FOOBAR"
	TestAccount = "EXB123456"
	NonVenue    = "NONVENUE"
)

func MakeStockfighterClient(t *testing.T) StockfighterClient {
	client, err := NewStockfighterClient()
	if err != nil {
		t.Fatalf("Failed to get client: %v\n", err)
	}
	return client
}

func TestHeartbeat(t *testing.T) {
	client := MakeStockfighterClient(t)
	if !client.Heartbeat() {
		t.Fatalf("Failed to detect a heartbeat.\n")
	}
}

func TestVenueHeartbeat(t *testing.T) {
	client := MakeStockfighterClient(t)
	if !client.VenueHeartbeat(TestVenue) {
		t.Fatalf("Failed to detect venue heartbeat.\n")
	}

	if client.VenueHeartbeat(NonVenue) {
		t.Fatalf("Got a heartbeat for a non-existent venue.\n")
	}
}

func TestStocks(t *testing.T) {
	client := MakeStockfighterClient(t)
	stocks, err := client.Stocks(TestVenue)
	if err != nil {
		t.Fatalf("Failed to get stocks: %v\n Stocks: %+v\n", err, stocks)
	}
	if len(stocks.Symbols) != 1 {
		t.Errorf("Expected one stock, got %v: %+v\n", len(stocks.Symbols), stocks.Symbols)
	}
	if stocks.Symbols[0].Symbol != "FOOBAR" {
		t.Errorf("Expected single stock to be FOOBAR, got %+v\n", stocks.Symbols[0])
	}

	stocks, err = client.Stocks(NonVenue)
	if err == nil {
		t.Fatalf("Expected error for non-existent venue: %+v\n", stocks)
	}
	if stocks.OK {
		t.Fatalf("Got OK for non-existent venue: %+v\n", stocks)
	}
	if stocks.Error == "" {
		t.Fatalf("Expected error for non-existent venue: %+v\n", stocks)
	}
}

func TestOrders(t *testing.T) {
	client := MakeStockfighterClient(t)
	orders, err := client.Orders(TestVenue, TestStock)
	if err != nil {
		t.Fatalf("Failed to get orders: %v\n Orders: %+v\n", err, orders)
	}
}

func TestPostOrders(t *testing.T) {
	client := MakeStockfighterClient(t)
	postOrder, err := client.PostOrder(PostOrderRequest{
		Venue:     TestVenue,
		Account:   TestAccount,
		Stock:     TestStock,
		Price:     100,
		Qty:       100,
		Direction: Buy,
		OrderType: Limit,
	})
	if err != nil {
		t.Fatalf("Error posting order: %v\n Order: %+v\n", err, postOrder)
	}
}

func TestQuote(t *testing.T) {
	client := MakeStockfighterClient(t)
	quote, err := client.Quote(TestVenue, TestStock)
	if err != nil {
		t.Fatalf("Error getting quote: %v\n Quote: %+v\n", err, quote)
	}
}

func TestOrder(t *testing.T) {
	client := MakeStockfighterClient(t)
	postOrder, err := client.PostOrder(PostOrderRequest{
		Venue:     TestVenue,
		Account:   TestAccount,
		Stock:     TestStock,
		Price:     100,
		Qty:       100,
		Direction: Buy,
		OrderType: Limit,
	})
	if err != nil {
		t.Fatalf("Error posting order: %v\n Order: %+v\n", err, postOrder)
	}

	orderStatus, err := client.Order(TestVenue, TestStock, postOrder.ID)
	if err != nil {
		t.Fatalf("Error getting status of order: %v Status: %+v\n", err, orderStatus)
	}
}

func TestCancel(t *testing.T) {
	client := MakeStockfighterClient(t)
	postOrder, err := client.PostOrder(PostOrderRequest{
		Venue:     TestVenue,
		Account:   TestAccount,
		Stock:     TestStock,
		Price:     100,
		Qty:       100,
		Direction: Buy,
		OrderType: Limit,
	})
	if err != nil {
		t.Fatalf("Error posting order: %v\n Order: %+v\n", err, postOrder)
	}

	order, err := client.Cancel(TestVenue, TestStock, postOrder.ID)
	if err != nil {
		t.Fatalf("Error cancelling order: %v\n Order: %+v\n", err, order)
	}
	if order.Open {
		t.Fatalf("Expected order to be cancelled! Order: %+v\n", order)
	}
}

func TestAccountOrders(t *testing.T) {
	client := MakeStockfighterClient(t)
	orders, err := client.AccountOrders(TestVenue, TestAccount)
	if err != nil {
		t.Fatalf("Error getting account orders: %v\n Order: %+v\n", err, orders)
	}
}

func TestAccountStockOrders(t *testing.T) {
	client := MakeStockfighterClient(t)
	orders, err := client.AccountStockOrders(TestVenue, TestAccount, TestStock)
	if err != nil {
		t.Fatalf("Error getting account orders: %v\n Order: %+v\n", err, orders)
	}
}
