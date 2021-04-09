package ping

import (
	"net"

	"github.com/miekg/dns"
	"github.com/secure-dns/service/plugin"
)

func init() {
	plugin.Register("ping", Ping)
}

var Ping = plugin.Plugin{
	Exec: exec,
}

func exec(req *dns.Msg) ([]dns.RR, uint8) {
	host := req.Question[0].Name
	qType := req.Question[0].Qtype

	if host != "ping.alxs.xyz." {
		return []dns.RR{}, plugin.Next
	}
	if qType == dns.TypeAAAA {
		r := new(dns.AAAA)
		r.Hdr = dns.RR_Header{Name: host, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 300}
		r.AAAA = net.ParseIP("2a03:4000:30:be8e::14:9493")
		return []dns.RR{r}, plugin.Stop
	}
	return []dns.RR{}, plugin.Stop
}
