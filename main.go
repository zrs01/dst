package main

import (
	"github.com/zrs01/dst/cmd/dst"
	"github.com/zrs01/dst/config"
)

var version = "development"

func main() {
	config.Version = version
	dst.Launch()
}
