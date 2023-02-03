package task

import (
	"LiveProxySpeedTest/internal/common"
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

func loadIPRanges(fileName string) (ips []*common.CustomIPAddr) {
	resp, err := req.C().R().Get("https://gh-proxy.com/https://raw.githubusercontent.com/sec-an/LiveProxySpeedTest/main/data/" + fileName + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	for _, ipInfo := range strings.Split(resp.String(), "\n") {
		info := strings.Split(strings.TrimSpace(ipInfo), "#")
		if len(info) > 1 {
			ips = append(ips, &common.CustomIPAddr{
				IPAddr: net.IPAddr{IP: net.ParseIP(info[0])},
				Loc:    info[1],
			})
			if err == io.EOF {
				return
			}

		}
	}
	return
}
