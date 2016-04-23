package stockfighter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	StockfighterAPIKeyEnvVar = "STOCKFIGHTER_API_KEY"
	StarfighterAuthHeader    = "X-Starfighter-Authorization"
	BaseURL                  = "https://api.stockfighter.io/ob/api"
)

type StockfighterClient struct {
	Client http.Client
	APIKey string
}

func NewStockfighterClient() (StockfighterClient, error) {
	client := StockfighterClient{
		APIKey: os.Getenv(StockfighterAPIKeyEnvVar),
	}
	if client.APIKey == "" {
		return client, fmt.Errorf("%v not set.", StockfighterAPIKeyEnvVar)
	}
	return client, nil
}

type StockfighterResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type StockfighterError interface {
	GetStockfighterError() error
}

func (s StockfighterResponse) GetStockfighterError() error {
	if !s.OK {
		return fmt.Errorf("Response not OK: %v", s.Error)
	}
	return nil
}

func (s StockfighterClient) Do(method string, url string, body interface{}, response StockfighterError) error {
	var bodyReader io.Reader
	bodyReader = nil
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("Failed to marshal body: %v\n", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	request, err := http.NewRequest(method, BaseURL+url, bodyReader)
	if err != nil {
		return fmt.Errorf("Failed to make request: %v", err)
	}
	request.Header.Add(StarfighterAuthHeader, s.APIKey)
	resp, err := s.Client.Do(request)
	if err != nil {
		return fmt.Errorf("Error getting URL: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("Error decoding response: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Non-200 status code: %v", resp.StatusCode)
	}

	return response.GetStockfighterError()
}

func (s StockfighterClient) Heartbeat() bool {
	stockfighterResponse := StockfighterResponse{}
	err := s.Do("GET", "/heartbeat", nil, &stockfighterResponse)
	if err != nil {
		return false
	}
	return true
}

func (s StockfighterClient) VenueHeartbeat(venue string) bool {
	stockfighterResponse := StockfighterResponse{}
	err := s.Do("GET", fmt.Sprintf("/venues/%v/heartbeat", venue), nil, &stockfighterResponse)
	if err != nil {
		return false
	}
	return true
}

type StockSymbol struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type StocksResponse struct {
	StockfighterResponse
	Symbols []StockSymbol `json:"symbols"`
}

func (s StockfighterClient) Stocks(venue string) (StocksResponse, error) {
	stocksResponse := StocksResponse{}
	err := s.Do("GET", fmt.Sprintf("/venues/%v/stocks", venue), nil, &stocksResponse)
	return stocksResponse, err
}

type Order struct {
	Price int  `json:"price"`
	Qty   int  `json:"qty"`
	IsBuy bool `json:"isBuy"`
}

type OrdersResponse struct {
	StockfighterResponse
	Venue     string    `json:"venue"`
	Symbol    string    `json:"symbol"`
	Bids      []Order   `json:"bids"`
	Asks      []Order   `json:"asks"`
	Timestamp time.Time `json:"ts"`
}

func (s StockfighterClient) Orders(venue string, stock string) (OrdersResponse, error) {
	ordersResponse := OrdersResponse{}
	err := s.Do("GET", fmt.Sprintf("/venues/%v/stocks/%v", venue, stock), nil, &ordersResponse)
	return ordersResponse, err
}

type Direction string

const (
	Buy  Direction = "buy"
	Sell Direction = "sell"
)

type OrderType string

const (
	Limit             OrderType = "limit"
	Market            OrderType = "market"
	FillOrKill        OrderType = "fill-or-kill"
	ImmediateOrCancel OrderType = "immediate-or-cancel"
)

type PostOrderRequest struct {
	Account   string    `json:"account"`
	Venue     string    `json:"venue"`
	Stock     string    `json:"stock"`
	Price     int       `json:"price"`
	Qty       int       `json:"qty"`
	Direction Direction `json:"direction"`
	OrderType OrderType `json:"orderType"`
}

type OrderStatusResponse struct {
	StockfighterResponse
	OrderStatus
}

type OrderStatus struct {
	Venue       string    `json:"venue"`
	Symbol      string    `json:"symbol"`
	Direction   Direction `json:"direction"`
	OriginalQty int       `json:"originalQty"`
	Price       int       `json:"price"`
	OrderType   OrderType `json:"orderType"`
	ID          int       `json:"id"`
	Account     string    `json:"account"`
	Timestamp   time.Time `json:"ts"`
	Fills       []Fill    `json:"fills"`
	TotalFilled int       `json:"totalFilled"`
	Open        bool      `json:"open"`
}

type Fill struct {
	Price     int       `json:"price"`
	Qty       int       `json:"qty"`
	Timestamp time.Time `json:"ts"`
}

// Note: /venues/:venue/stocks/:stock seems to allow POST as an internal endpoint, but needs an internal key.
func (s StockfighterClient) PostOrder(request PostOrderRequest) (OrderStatusResponse, error) {
	orderStatus := OrderStatusResponse{}
	err := s.Do("POST", fmt.Sprintf("/venues/%v/stocks/%v/orders", request.Venue, request.Stock), request, &orderStatus)
	return orderStatus, err
}

type QuoteResponse struct {
	StockfighterResponse
	Venue     string    `json:"venue"`
	Symbol    string    `json:"symbol"`
	Bid       int       `json:"bid"`
	Ask       int       `json:"ask"`
	BidSize   int       `json:"bidSize"`
	AskSize   int       `json:"askSize"`
	BidDepth  int       `json:"bidDepth"`
	AskDepth  int       `json:"askDepth"`
	Last      int       `json:"last"`
	LastSize  int       `json:"lastSize"`
	LastTrade time.Time `json:"lastTrade"`
	Timestamp time.Time `json:"quoteTime"`
}

func (s StockfighterClient) Quote(venue string, stock string) (QuoteResponse, error) {
	quoteResponse := QuoteResponse{}
	err := s.Do("GET", fmt.Sprintf("/venues/%v/stocks/%v/quote", venue, stock), nil, &quoteResponse)
	return quoteResponse, err
}

type OrderResponse struct {
	StockfighterResponse
}

func (s StockfighterClient) Order(venue string, stock string, order int) (OrderStatusResponse, error) {
	orderStatus := OrderStatusResponse{}
	err := s.Do("GET", fmt.Sprintf("/venues/%v/stocks/%v/orders/%v", venue, stock, order), nil, &orderStatus)
	return orderStatus, err
}

func (s StockfighterClient) Cancel(venue string, stock string, order int) (OrderStatusResponse, error) {
	orderStatus := OrderStatusResponse{}
	err := s.Do("DELETE", fmt.Sprintf("/venues/%v/stocks/%v/orders/%v", venue, stock, order), nil, &orderStatus)
	return orderStatus, err
}

type AccountOrdersResponse struct {
	StockfighterResponse
	Venue  string        `json:"venue"`
	Orders []OrderStatus `json:"orders"`
}

func (s StockfighterClient) AccountOrders(venue string, account string) (AccountOrdersResponse, error) {
	accountOrdersResponse := AccountOrdersResponse{}
	err := s.Do("GET", fmt.Sprintf("/venues/%v/accounts/%v/orders", venue, account), nil, &accountOrdersResponse)
	return accountOrdersResponse, err
}

func (s StockfighterClient) AccountStockOrders(venue string, account string, stock string) (AccountOrdersResponse, error) {
	accountOrdersResponse := AccountOrdersResponse{}
	err := s.Do("GET", fmt.Sprintf("/venues/%v/accounts/%v/stocks/%v/orders", venue, account, stock), nil, &accountOrdersResponse)
	return accountOrdersResponse, err
}
