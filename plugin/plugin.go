package plugin

import (
	"log"

	"github.com/miekg/dns"
)

//Plugin - default plugin type
type Plugin struct {
	//Exec is the only required function of a plugin. It handels a user request.
	Exec func(req *dns.Msg) ([]dns.RR, uint8)
	//Cron can be used to update the plugin data
	Cron func()
}

//PluginList - list of all existing plugins
type PluginList map[string]Plugin

//Plugins - list of all available plugins
var Plugins = PluginList{}

const (
	Next uint8 = 0
	Stop uint8 = 1
)

//Register is a function to register new plugins.
func Register(name string, plugin Plugin) {

	if Plugins[name].Exec != nil {
		log.Fatalf("plugin \"%s\" already exists", name)
		return
	}

	Plugins[name] = plugin
}
