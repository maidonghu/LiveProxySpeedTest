package task

import (
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/imroc/req/v3"
)

func InitRandSeed() {
	rand.Seed(time.Now().UnixNano())
}

func isIPv4(ip string) bool {
	return strings.Contains(ip, ".")
}

func loadIPRanges(fileName string) (ips []*net.IPAddr) {
	resp, err := req.C().R().Get("https://gh-proxy.com/https://raw.githubusercontent." +
		"com/sec-an/LiveProxySpeedTest/main/data/" + fileName + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	for _, ip := range strings.Fields(resp.String()) {
		ips = append(ips, &net.IPAddr{IP: net.ParseIP(strings.TrimSpace(ip))})
		if err == io.EOF {
			return
		}

	}
	return
}
