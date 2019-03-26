package integer

import "testing"

func TestProcessEvent(t *testing.T) {
	type args struct {
		replacement interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"smoke test",
			args{
				replacement: []int{1, 2},
			},
			"1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessEvent(tt.args.replacement); got != tt.want {
				t.Errorf("ProcessEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
