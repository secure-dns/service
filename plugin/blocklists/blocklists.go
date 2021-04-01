package blocklists

import (
	"bufio"
	"bytes"
	"net/http"
	"sort"

	cplugin "github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"github.com/secure-dns/service/plugin"
)

func GetPlugin(path string) plugin.Plugin {
	hosts := readData(path)
	return plugin.Plugin{
		Exec: func(req *dns.Msg) ([]dns.RR, uint8) {
			return exec(req, hosts)
		},
		Cron: func() {
			hosts = readData(path)
		},
	}
}

func exec(req *dns.Msg, hosts []string) ([]dns.RR, uint8) {
	host := req.Question[0].Name
	if !isOnList(hosts, host) {
		return []dns.RR{}, plugin.Next
	}

	return []dns.RR{}, plugin.Stop
}

func isOnList(hosts []string, host string) bool {
	if i := sort.SearchStrings(hosts, host); i < len(hosts) && hosts[i] == host {
		return true
	}
	return false
}

func readData(path string) []string {
	//data, err := ioutil.ReadFile(path)
	resp, err := http.Get(path)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
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
