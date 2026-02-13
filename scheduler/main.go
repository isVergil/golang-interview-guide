package main

import (
	"fmt"
	"time"
)

func main() {
	triggerTime := time.Unix(0, 0).Format("2006-01-02 15:04:05")
	fmt.Println(triggerTime)
}
