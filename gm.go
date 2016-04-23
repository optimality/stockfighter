package stockfighter

import "fmt"

const (
	GMURL = "/gm"
)

type InstanceResponse struct {
	StockfighterResponse
	Account              string            `json:"account"`
	Balances             map[string]int    `json:"balances"`
	InstanceID           int               `json:"instanceId"`
	Instructions         map[string]string `json:"instructions"`
	SecondsPerTradingDay int               `json:"secondsPerTradingDay"`
	Tickers              []string          `json:"tickers"`
	Venues               []string          `json:"venues"`
}

func (s StockfighterClient) StartLevel(level string) (InstanceResponse, error) {
	instanceResponse := InstanceResponse{}
	err := s.Do("POST", fmt.Sprintf(GMURL+"/levels/%v", level), nil, &instanceResponse)
	return instanceResponse, err
}

func (s StockfighterClient) RestartInstance(instance int) (InstanceResponse, error) {
	instanceResponse := InstanceResponse{}
	err := s.Do("POST", fmt.Sprintf(GMURL+"/instances/%v/restart", instance), nil, &instanceResponse)
	return instanceResponse, err
}

func (s StockfighterClient) ResumeInstance(instance int) (InstanceResponse, error) {
	instanceResponse := InstanceResponse{}
	err := s.Do("POST", fmt.Sprintf(GMURL+"/instances/%v/resume", instance), nil, &instanceResponse)
	return instanceResponse, err
}

func (s StockfighterClient) StopInstance(instance int) (StockfighterResponse, error) {
	stockfighterResponse := StockfighterResponse{}
	err := s.Do("POST", fmt.Sprintf(GMURL+"/instances/%v/stop", instance), nil, &stockfighterResponse)
	return stockfighterResponse, err
}
