package core

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/coredns/coredns/plugin/pkg/doh"
	"github.com/miekg/dns"
	"github.com/secure-dns/service/handler"
)

//runDoH starts the http server
func runDoH(addr string) {
	http.HandleFunc("/dns-query/", dohHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}

//dohHandler handels http requests
func dohHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := r.URL.Parse(r.RequestURI)
	msg, err := doh.RequestToMsg(r)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	plugins := strings.Split(strings.TrimPrefix(u.Path, "/dns-query/"), ":")

	resp := handler.HandlePlugins(msg, plugins)

	err = sendDohResult(w, resp)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

}

//sendErrorResponse sends a bad request back to the client.
func sendErrorResponse(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

//sendDohResult is sending valid dns responses.
func sendDohResult(w http.ResponseWriter, resp *dns.Msg) error {
	buf, err := resp.Pack()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", doh.MimeType)
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%f", 2.5))
	w.Header().Set("Content-Length", strconv.Itoa(len(buf)))
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
	return nil
}
