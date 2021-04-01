package handler

import (
	"sort"

	"github.com/miekg/dns"
	"github.com/secure-dns/service/plugin"
)

//HandlePlugins is called on each request and fills the gap between the client and your plugins.
func HandlePlugins(req *dns.Msg, plugins []string) *dns.Msg {

	if i := sort.SearchStrings(plugins, "fwd"); i >= len(plugins) || plugins[i] != "fwd" {
		plugins = append(plugins, "fwd")
	}

	msg := new(dns.Msg)
	msg.SetReply(req)
	msg.Authoritative = true

	answers := []dns.RR{}
	//run plugins

	for _, pluginName := range plugins {
		currentPlugin := plugin.Plugins[pluginName]
		if currentPlugin.Exec == nil {
			continue
		}
		res, code := currentPlugin.Exec(req)
		if len(res) != 0 {
			answers = append(answers, res...)
			break
		}
		if code == plugin.Stop {
			break
		}
	}

	//bind results to msg
	msg.Answer = answers

	return msg
}
