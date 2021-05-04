package forward

import (
	"net"

	"github.com/miekg/dns"
	"github.com/secure-dns/service/plugin"
)

func init() {
	plugin.Register("fwd", Forward)
}

var Forward = plugin.Plugin{
	Exec: exec,
}

func exec(req *dns.Msg) ([]dns.RR, uint8) {
	c := new(dns.Client)
	res, _, err := c.Exchange(req, net.JoinHostPort("1.1.1.1", "53"))
	if err != nil {
		return []dns.RR{}, plugin.Next
	}
	return res.Answer, plugin.Stop
}
