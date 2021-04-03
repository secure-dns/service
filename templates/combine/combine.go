package combine

import (
	"github.com/miekg/dns"
	"github.com/secure-dns/service/handler"
	"github.com/secure-dns/service/plugin"
)

func GetPlugin(plugins []string) plugin.Plugin {
	return plugin.Plugin{
		Exec: func(req *dns.Msg) ([]dns.RR, uint8) {
			res := handler.RunPluginChain(plugins, req)
			if len(res) == 0 {
				return []dns.RR{}, plugin.Next
			}

			return res, plugin.Stop
		},
	}
}
