package webhook

import "testing"

// TestParse garante que payloads válidos são aceitos e inválidos rejeitados.
func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantRef   string
		wantText  string
	}{
		{
			name:     "válido completo",
			input:    `{"customer_ref": "u-1", "text": "oi"}`,
			wantRef:  "u-1",
			wantText: "oi",
		},
		{
			name:    "faltando text",
			input:   `{"customer_ref": "u-1"}`,
			wantErr: true,
		},
		{
			name:    "faltando customer_ref",
			input:   `{"text": "oi"}`,
			wantErr: true,
		},
		{
			name:    "JSON inválido",
			input:   `{`,
			wantErr: true,
		},
		{
			name:    "objeto vazio",
			input:   `{}`,
			wantErr: true,
		},
		{
			name:    "outros campos ignorados",
			input:   `{"customer_ref": "u-1", "text": "oi", "timestamp": 123}`,
			wantRef:  "u-1",
			wantText: "oi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse(%q) err = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got.CustomerRef != tt.wantRef {
				t.Fatalf("CustomerRef = %q, want %q", got.CustomerRef, tt.wantRef)
			}
			if got.Text != tt.wantText {
				t.Fatalf("Text = %q, want %q", got.Text, tt.wantText)
			}
		})
	}
}
