package marusia

import "encoding/json"

type ResponseBody struct {
	Response `json:"response"`
	Session  `json:"session"`
	Version  string `json:"version"`
}

type Response struct {
	Text       string   `json:"text"`
	TTS        string   `json:"tts,omitempty"`
	Buttons    []Button `json:"buttons,omitempty"`
	EndSession bool     `json:"end_session"`
}

type Button struct {
	Title   string          `json:"title"`
	Payload json.RawMessage `json:"payload,omitempty"`
	Url     string          `json:"url,omitempty"`
}
