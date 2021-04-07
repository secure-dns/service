package core

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/miekg/dns"
	"github.com/secure-dns/service/handler"
)

var ServerName string

type dnsHandler struct{}

func (this *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if !strings.HasSuffix(ServerName, os.Getenv("TLS_DOMAIN_SUFFIX")) {
		w.Close()
		return
	}
	plugins := strings.Split(strings.TrimSuffix(ServerName, os.Getenv("TLS_DOMAIN_SUFFIX")), "-")

	resp := handler.HandlePlugins(r, plugins)
	buf, err := resp.Pack()
	if err != nil {
		return
	}

	w.Write(buf)
}

//runDoT starts a dns over tls server.
func runDoT(addr string) {
	cer, err := tls.LoadX509KeyPair(os.Getenv("TLS_CERT_CHAIN"), os.Getenv("TLS_CERT_KEY"))
	if err != nil {
		log.Println(err)
		return
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cer},
		GetConfigForClient: func(argHello *tls.ClientHelloInfo) (*tls.Config, error) {
			hello := new(tls.ClientHelloInfo)
			*hello = *argHello
			ServerName = hello.ServerName
			return nil, nil
		}}

	srv := &dns.Server{Addr: addr, Net: "tcp-tls", TLSConfig: config}
	srv.Handler = &dnsHandler{}
	fmt.Printf("DoT Server listening on: %s\n", addr)
	log.Fatal(srv.ListenAndServe())
}
