package exception

import (
	"errors"
	"testing"
)

func TestIsDataNotFound(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case1",
			args: args{
				err: Custom(DataNotFound, "数据不存在"),
			},
			want: true,
		},
		{
			name: "case2",
			args: args{
				err: errors.New("asdf"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDataNotFound(tt.args.err); got != tt.want {
				t.Errorf("IsDataNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}
