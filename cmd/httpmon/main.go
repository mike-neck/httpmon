package main

import (
	"fmt"
	"github.com/mike-neck/httpmon"
)

func main() {
	t, _ := httpmon.TimeOutFromString("5s")
	fmt.Println(t)
}
