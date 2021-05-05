package handler

import (
	"github.com/miekg/dns"
	"github.com/secure-dns/service/plugin"
)

//HandlePlugins is called on each request and fills the gap between the client and your plugins.
func HandlePlugins(req *dns.Msg, plugins []string) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetReply(req)
	msg.Authoritative = true

	//run plugins
	answers := RunPluginChain(plugins, req)

	//bind results to msg
	msg.Answer = answers

	return msg
}

//RunPluginChain - runs a list of plugins
func RunPluginChain(pluginList []string, req *dns.Msg) []dns.RR {
	for _, pluginName := range pluginList {
		currentPlugin := plugin.Plugins[pluginName]
		if currentPlugin.Exec == nil {
			continue
		}
		res, code := currentPlugin.Exec(req)
		if code == plugin.Stop {
			return res
		}
	}
	return []dns.RR{}
}
