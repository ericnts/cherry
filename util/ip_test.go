package util

import "testing"

func Test_localIP(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "case1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LocalIP(); got != tt.want {
				t.Errorf("LocalIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
