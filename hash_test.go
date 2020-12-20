package filedriller

import (
	"reflect"
	"testing"
)

func TestHashit(t *testing.T) {
	type args struct {
		inFile  string
		hashalg string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"md5", args{"pronom.sig", "md5"}, []byte{245, 231, 48, 19, 74, 17, 59, 216, 138, 151, 23, 150, 36, 136, 10, 1}},
		{"sha1", args{"pronom.sig", "sha1"}, []byte{150, 107, 206, 220, 202, 56, 151, 113, 214, 4, 152, 88, 118, 9, 51, 93, 0, 36, 120, 248}},
		{"sha256", args{"pronom.sig", "sha256"}, []byte{252, 234, 8, 28, 247, 21, 104, 147, 168, 14, 140, 99, 49, 124, 203, 199, 96, 177, 173, 213, 82, 240, 32, 186, 16, 45, 194, 46, 120, 71, 232, 129}},
		{"sha512", args{"pronom.sig", "sha512"}, []byte{50, 194, 92, 137, 35, 66, 100, 180, 139, 169, 38, 197, 238, 156, 208, 144, 103, 119, 50, 90, 243, 7, 114, 163, 237, 18, 74, 229, 148, 101, 245, 75, 95, 228, 195, 64, 152, 147, 218, 197, 57, 223, 216, 181, 246, 253, 164, 43, 216, 69, 5, 163, 67, 234, 117, 152, 216, 177, 198, 178, 118, 40, 155, 130}},
		{"blake2-512", args{"pronom.sig", "blake2b-512"}, []byte{112, 94, 50, 70, 33, 102, 235, 121, 132, 172, 11, 160, 39, 245, 227, 219, 23, 181, 252, 35, 72, 102, 157, 201, 215, 22, 223, 79, 5, 69, 95, 96, 224, 207, 190, 167, 76, 194, 25, 169, 57, 39, 203, 133, 132, 152, 216, 207, 92, 65, 212, 89, 103, 1, 54, 49, 94, 93, 22, 239, 162, 114, 9, 180}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hashit(tt.args.inFile, tt.args.hashalg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hashit() = %v, want %v", got, tt.want)
			}
		})
	}
}
