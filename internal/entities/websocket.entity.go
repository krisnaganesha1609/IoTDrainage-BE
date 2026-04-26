package entities

type WebsocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}
