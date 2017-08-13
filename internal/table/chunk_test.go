package table

import (
	"fmt"
	"testing"
)

func TestChunkQuery(t *testing.T) {
	chunk := chunk{
		data: [][]int{
			{
				1, 2, 3, 4, 5,
			},
			{
				11, 12, 13, 14, 15,
			},
			{
				21, 22, 23, 24, 25,
			},
		},
	}

	next := chunk.query(2, 2)

	for next.Next() {
		fmt.Println(next.Value())
	}

}
