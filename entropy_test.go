package filedriller

import "testing"

func Test_entropy(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name        string
		args        args
		wantEntropy float64
		wantErr     bool
	}{
		{"Calculate entropy", args{"testdata/1200px-GPLv3_Logo.svg.png"}, 7.96, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntropy, err := entropy(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("entropy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEntropy != tt.wantEntropy {
				t.Errorf("entropy() gotEntropy = %v, want %v", gotEntropy, tt.wantEntropy)
			}
		})
	}
}
