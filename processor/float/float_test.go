package float

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
			name: "smoke test",
			args: args{
				replacement: []float64{1.0, 2.0, 3},
			},
			want: "1.605",
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
