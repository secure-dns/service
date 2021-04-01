package hostfile

import (
	"github.com/miekg/dns"
	"github.com/secure-dns/service/plugin"
)

const ttl = uint32(600)

func GetPlugin(path string, action uint8) plugin.Plugin {
	hosts, _ := parse(path, action)

	return plugin.Plugin{
		Exec: func(req *dns.Msg) []dns.RR {
			qName := req.Question[0].Name
			qType := req.Question[0].Qclass

			answers := []dns.RR{}

			switch qType {
			case dns.TypePTR:
				names := LookupStaticAddr(*hosts, qName)
				if len(names) == 0 {
					return answers
				}
				answers = ptr(qName, ttl, names)
			case dns.TypeA:
				ips := LookupStaticHostV4(*hosts, qName)
				answers = a(qName, ttl, ips)
			case dns.TypeAAAA:
				ips := LookupStaticHostV6(*hosts, qName)
				answers = aaaa(qName, ttl, ips)
			default:
				answers = []dns.RR{}
			}
			return answers
		},
		Cron: func() {
			hosts, _ = parse(path, action)
		},
	}
}

const (
	//TypeWeb - path is a url string
	TypeWeb uint8 = 0
	//TypeFile - path is a local file path
	TypeFile uint8 = 1
)
