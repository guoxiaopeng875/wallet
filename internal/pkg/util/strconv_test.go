package util

import "testing"

func TestStringToUint(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		wantErr bool
	}{
		{"success", args{"123"}, 123, false},
		{"fail", args{"abc"}, 0, true},
		{"fail", args{"-123"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToUint(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringToUint() got = %v, want %v", got, tt.want)
			}
		})
	}
}
