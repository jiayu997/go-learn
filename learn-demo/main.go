package main

import (
	"fmt"
	"os"
)

func main() {
	test := os.Getenv("UPDATE_POLICY")
	if test != "CREATE" && test != "UPDATE" {
		test = "failed"
	} else {
		test = "ok"
	}
	fmt.Println(test)
}
