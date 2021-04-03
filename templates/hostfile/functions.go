package hostfile

//source: https://github.com/coredns/coredns

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type Hosts struct {
	name4 map[string][]net.IP
	name6 map[string][]net.IP
	addr  map[string][]string
}

func newHosts() *Hosts {
	return &Hosts{
		name4: make(map[string][]net.IP),
		name6: make(map[string][]net.IP),
		addr:  make(map[string][]string),
	}
}

// Parse reads the hostsfile and populates the byName and addr maps.
func parse(path string, action uint8) (*Hosts, error) {
	var r io.Reader
	switch action {
	case TypeFile:
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return newHosts(), err
		}
		r = bytes.NewReader(b)
	default:
		resp, err := http.Get(path)
		if err != nil {
			return newHosts(), err
		}
		defer resp.Body.Close()
		r = resp.Body
	}

	hmap := newHosts()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if i := bytes.Index(line, []byte{'#'}); i >= 0 {
			//Discard comments.
			line = line[0:i]
		}
		f := bytes.Fields(line)
		if len(f) < 2 {
			continue
		}
		addr := net.ParseIP(string(f[0]))
		if addr == nil {
			continue
		}

		family := 2
		if addr.To4() != nil {
			family = 1
		}

		for i := 1; i < len(f); i++ {
			name := plugin.Name(string(f[i])).Normalize()
			switch family {
			case 1:
				hmap.name4[name] = append(hmap.name4[name], addr)
			case 2:
				hmap.name6[name] = append(hmap.name6[name], addr)
			default:
				continue
			}
			hmap.addr[addr.String()] = append(hmap.addr[addr.String()], name)
		}

	}
	return hmap, nil
}

// lookupStaticHost looks up the IP addresses for the given host from the hosts file.
func lookupStaticHost(m map[string][]net.IP, host string) []net.IP {
	if len(m) == 0 {
		return nil
	}
	ips, ok := m[host]
	if !ok {
		return nil
	}
	ipsCp := make([]net.IP, len(ips))
	copy(ipsCp, ips)
	return ipsCp
}

// LookupStaticHostV4 looks up the IPv4 addresses for the given host from the hosts file.
func LookupStaticHostV4(m Hosts, host string) []net.IP {
	host = strings.ToLower(host)
	return lookupStaticHost(m.name4, host)
}

// LookupStaticHostV6 looks up the IPv6 addresses for the given host from the hosts file.
func LookupStaticHostV6(m Hosts, host string) []net.IP {
	host = strings.ToLower(host)
	return lookupStaticHost(m.name6, host)
}

// parseIP calls discards any v6 zone info, before calling net.ParseIP.
func parseIP(addr string) net.IP {
	if i := strings.Index(addr, "%"); i >= 0 {
		// discard ipv6 zone
		addr = addr[0:i]
	}

	return net.ParseIP(addr)
}

// LookupStaticAddr looks up the hosts for the given address from the hosts file.
func LookupStaticAddr(m Hosts, addr string) []string {
	addr = parseIP(addr).String()
	if addr == "" {
		return nil
	}

	hosts1 := m.addr[addr]

	if len(hosts1) == 0 {
		return nil
	}

	hostsCp := make([]string, len(hosts1))
	copy(hostsCp, hosts1)
	return hostsCp
}

// a takes a slice of net.IPs and returns a slice of A RRs.
func a(zone string, ttl uint32, ips []net.IP) []dns.RR {
	answers := make([]dns.RR, len(ips))
	for i, ip := range ips {
		r := new(dns.A)
		r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl}
		r.A = ip
		answers[i] = r
	}
	return answers
}

// aaaa takes a slice of net.IPs and returns a slice of AAAA RRs.
func aaaa(zone string, ttl uint32, ips []net.IP) []dns.RR {
	answers := make([]dns.RR, len(ips))
	for i, ip := range ips {
		r := new(dns.AAAA)
		r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: ttl}
		r.AAAA = ip
		answers[i] = r
	}
	return answers
}

// ptr takes a slice of host names and filters out the ones that aren't in Origins, if specified, and returns a slice of PTR RRs.
func ptr(zone string, ttl uint32, names []string) []dns.RR {
	answers := make([]dns.RR, len(names))
	for i, n := range names {
		r := new(dns.PTR)
		r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypePTR, Class: dns.ClassINET, Ttl: ttl}
		r.Ptr = dns.Fqdn(n)
		answers[i] = r
	}
	return answers
}
