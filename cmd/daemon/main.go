package main

import (
	"github.com/secure-dns/service/core"
	_ "github.com/secure-dns/service/core-plugins"
)

func main() {
	core.Run(":8080")
}
