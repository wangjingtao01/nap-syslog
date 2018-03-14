package dispatch

type DispatcherConfig struct {
	Dispatcher []DispatcherInstance `json:"dispatcher"`
}

type DispatcherInstance struct {
	Type         string             `json:"type"`
	URI          string             `json:"uri,omitempty"`
	Exchange     string             `json:"exchange,omitempty"`
	ExchangeType string             `json:"exchangeType,omitempty"`
	Triggers     map[string]Trigger `json:"triggers,omitempty"`
}

type Trigger struct {
	Queue      string `json:"queue"`
	RoutingKey string `json:"routingKey"`
}
