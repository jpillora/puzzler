package p0001

import (
	"reflect"
	"testing"
)

func Test_twoSum(t *testing.T) {
	type args struct {
		nums   []int
		target int
	}
	tests := []struct {
		name   string
		input  args
		output []int
	}{
		{
			name:  "test1",
			input: args{[]int{2, 7, 11, 15}, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := twoSum(tt.input.nums, tt.input.target); !reflect.DeepEqual(got, tt.output) {
				t.Errorf("twoSum() = %v, want %v", got, tt.output)
			}
		})
	}
}
