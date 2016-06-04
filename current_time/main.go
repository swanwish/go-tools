package main

import (
	"time"

	"github.com/swanwish/go-common/logs"
)

func main() {
	logs.Debugf("The current time unix is %d", time.Now().Unix())
}
