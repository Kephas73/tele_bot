package model

type RawData struct {
	Text   string      `json:"text,omitempty"`
	Object interface{} `json:"object,omitempty"`
}