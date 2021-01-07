package filedriller

import "testing"

func TestCreateUUID(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"Test UUID length", 36},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(CreateUUID()); got != tt.want {
				t.Errorf("CreateUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
