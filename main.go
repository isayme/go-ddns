package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/isayme/go-ddns/src/conf"
	"github.com/isayme/go-ddns/src/dnspod"
	"github.com/isayme/go-ddns/src/ip"
	"github.com/isayme/go-ddns/src/util"
	logger "github.com/isayme/go-logger"
)

var showVersion = flag.Bool("v", false, "show version")

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("name: %s\nversion: %s\nbuildTime: %s\ngitSHA1: %s\n", util.Name, util.Version, util.BuildTime, util.GitSHA1)
		os.Exit(0)
	}

	cfg := conf.Get()

	if cfg.Logger.Level != "" {
		logger.SetLevel(cfg.Logger.Level)
	}

	dnspod := dnspod.NewDNSPod(cfg.DNSPod)
	for {
		ip, err := ip.Get()
		if err != nil {
			logger.Errorf("get ip fail: %v", err)
			dnspod.Sleep()
			continue
		}

		err = dnspod.UpdateIP(ip)
		if err != nil {
			logger.Errorf("update ip fail: %v", err)
			dnspod.Sleep()
			continue
		}

		logger.Infow("update ip success", "ip", ip)
		dnspod.Sleep()
	}
}
