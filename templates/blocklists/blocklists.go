package blocklists

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	cplugin "github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"github.com/secure-dns/service/plugin"
)

const (
	//TypeWeb - path is a url string
	TypeWeb uint8 = 0
	//TypeFile - path is a local file path
	TypeFile uint8 = 1
)

func GetPlugin(path string, action uint8) plugin.Plugin {
	hosts := readData(path, action)
	return plugin.Plugin{
		Exec: func(req *dns.Msg) ([]dns.RR, uint8) {
			return exec(req, hosts)
		},
		Cron: func() {
			hosts = readData(path, action)
		},
	}
}

func exec(req *dns.Msg, hosts []string) ([]dns.RR, uint8) {
	host := req.Question[0].Name
	if !isOnList(hosts, host) {
		return []dns.RR{}, plugin.Next
	}

	switch req.Question[0].Qtype {
	case dns.TypeA:
		r := new(dns.A)
		r.Hdr = dns.RR_Header{Name: host, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 600}
		r.A = net.ParseIP("0.0.0.0").To4()
		return []dns.RR{r}, plugin.Stop
	case dns.TypeAAAA:
		r := new(dns.AAAA)
		r.Hdr = dns.RR_Header{Name: host, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 600}
		r.AAAA = net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000").To16()
		return []dns.RR{r}, plugin.Stop
	default:
		return []dns.RR{}, plugin.Stop
	}
}

func isOnList(hosts []string, host string) bool {
	/*if i := sort.SearchStrings(hosts, host); i < len(hosts) && hosts[i] == host {
		return true
	}
	return false*/
	for _, h := range hosts {
		if h == host {
			return true
		}
	}
	return false
}

func readData(path string, action uint8) []string {
	var r io.Reader
	switch action {
	case TypeFile:
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return []string{}
		}
		r = bytes.NewReader(b)
	default:
		resp, err := http.Get(path)
		if err != nil {
			return []string{}
		}
		defer resp.Body.Close()
		r = resp.Body
	}
	scanner := bufio.NewScanner(r)
	data := []string{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if i := bytes.Index(line, []byte{'#'}); i >= 0 {
			//Discard comments.
			line = line[0:i]
		}
		f := bytes.Fields(line)
		if len(f) != 1 {
			continue
		}
		data = append(data, cplugin.Name(string(f[0])).Normalize())
	}
	return data
}
