package bootstrapclient

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"strings"
)

func StaticDnsTransport() *http.Transport {
	return StaticDnsTransportWith(map[string][]string{
		"dns10.quad9.net":    {"149.112.112.10", "9.9.9.10"},
		"dns.google.com":     {"8.8.8.8", "8.8.4.4"},
		"cloudflare-dns.com": {"104.16.249.249", "104.16.248.249"},
	})
}

func StaticDnsTransportWith(resolved map[string][]string) *http.Transport {
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host := addr
			suffix := ""
			p := strings.LastIndex(addr, ":")
			if p >= 0 {
				host = addr[:p]
				suffix = addr[p:]
			}
			found, ok := resolved[host]
			if !ok {
				return net.Dial(network, addr)
			}
			randIndex := rand.Intn(len(found))
			resolvedAddr := found[randIndex] + suffix
			return net.Dial(network, resolvedAddr)
		},
	}
}
