package marusia

import "encoding/json"

// RequestBody structure to be sent as POST
type RequestBody struct {
	Meta    `json:"meta"`
	Request `json:"request"`
	Session `json:"session"`
	Version string `json:"version"`
}

type Request struct {
	Command           string          `json:"command"`
	OriginalUtterance string          `json:"original_utterance"`
	Type              string          `json:"type"`
	Payload           json.RawMessage `json:"payload"`
	NLU               NLU             `json:"nlu"`
}

type NLU struct {
	Tokens []string `json:"tokens"`
}
