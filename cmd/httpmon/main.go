package main

import (
	"fmt"
	"github.com/mike-neck/httpmon"
)

func main() {
	m := httpmon.Run()
	fmt.Println(m)
}
