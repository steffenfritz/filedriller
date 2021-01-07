package filedriller

import "testing"

func TestDownloadPronom(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Download signature file", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DownloadPronom(); (err != nil) != tt.wantErr {
				t.Errorf("DownloadPronom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
