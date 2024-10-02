package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_chunkSlice(t *testing.T) {
	type args struct {
		s    []int
		size int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "even divide",
			args: args{
				s:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				size: 3,
			},
			want: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
		},
		{
			name: "uneven divide",
			args: args{
				s:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				size: 8,
			},
			want: [][]int{
				{1, 2, 3, 4, 5, 6, 7, 8},
				{9},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := chunkSlice(tt.args.s, tt.args.size)
			assert.Equal(t, tt.want, got)
		})
	}
}
