package forward

import (
	"crypto/tls"
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

func exec(req *dns.Msg) []dns.RR {
	c := new(dns.Client)
	c.Net = "tcp-tls"
	c.TLSConfig = new(tls.Config)
	c.TLSConfig.ServerName = "cloudflare-dns.com"
	res, _, err := c.Exchange(req, net.JoinHostPort("1.1.1.1", "853"))
	if err != nil {
		return []dns.RR{}
	}
	return res.Answer
}
