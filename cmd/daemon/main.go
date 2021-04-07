package main

import (
	"github.com/secure-dns/service/core"
	_ "github.com/secure-dns/service/core-plugins"
	"github.com/secure-dns/service/env"
)

func main() {
	env.Load()
	core.Run(core.CoreConfig{DoH: true, DoT: true})
}
