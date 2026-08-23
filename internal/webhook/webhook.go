package webhook

import (
	"encoding/json"
	"fmt"
)

type Payload struct {
	CustomerRef string `json:"customer_ref"`
	Text        string `json:"text"`
}

func Parse(data []byte) (Payload, error) {
	var p Payload
	if err := json.Unmarshal(data, &p); err != nil {
		return p, fmt.Errorf("webhook: JSON inválido: %w", err)
	}
	if p.CustomerRef == "" || p.Text == "" {
		return p, fmt.Errorf("webhook: customer_ref e text são obrigatórios")
	}
	return p, nil
}
