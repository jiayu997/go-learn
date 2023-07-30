package main

import (
	"fmt"
)

func main() {
LABEL1:
	for k := 0; k < 3; k++ {
		for i := 0; i <= 5; i++ {
			for j := 0; j <= 5; j++ {
				if j == 4 {
					continue LABEL1
				}
				fmt.Printf("k is: %d,i is: %d, and j is: %d\n", k, i, j)
			}
		}
	}
}
