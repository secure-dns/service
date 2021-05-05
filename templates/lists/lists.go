package lists

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/secure-dns/service/plugin"
)

func GetPlugin(path string) plugin.Plugin {
	p := plgn{listPath: path}
	return plugin.Plugin{
		Exec: p.exec,
	}
}

type plgn struct {
	listPath string
}

func (p plgn) exec(req *dns.Msg) ([]dns.RR, uint8) {
	host := strings.TrimSuffix(req.Question[0].Name, ".")

	switch req.Question[0].Qtype {
	case dns.TypeA:
		if !p.check(host) {
			break
		}

		r := new(dns.A)
		r.Hdr = dns.RR_Header{Name: req.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 600}
		r.A = net.ParseIP("0.0.0.0").To4()

		return []dns.RR{r}, plugin.Stop
	case dns.TypeAAAA:
		if !p.check(host) {
			break
		}

		r := new(dns.AAAA)
		r.Hdr = dns.RR_Header{Name: req.Question[0].Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 600}
		r.AAAA = net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000").To16()

		return []dns.RR{r}, plugin.Stop
	}
	return []dns.RR{}, plugin.Next
}

func (p plgn) check(host string) bool {
	filename := host[0:2]
	if len(strings.Split(host, ".")[0]) == 1 {
		filename = "_" + host[0:1]
	}
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.txt", p.listPath, filename))
	if err != nil {
		return false
	}
	i := bytes.Index(dat, []byte(host))
	if i != -1 {
		var x int
		var y int
		for x = i; x > 0; x-- {
			if dat[x] == byte('\n') {
				break
			}
		}
		for y = i; y < len(dat); y++ {
			if dat[y] == byte('\n') {
				break
			}
		}
		//fmt.Println(string(dat[x : y+1]))
		return true
	}
	return false
}
