package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/utils"
	"github.com/networkservicemesh/networkservicemesh/utils/caddyfile"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
	_ "github.com/coredns/coredns/plugin/bind"
	_ "github.com/coredns/coredns/plugin/hosts"
	_ "github.com/coredns/coredns/plugin/log"
	_ "github.com/coredns/coredns/plugin/reload"

	"github.com/networkservicemesh/networkservicemesh/k8s/cmd/nsm-coredns/env"
	_ "github.com/networkservicemesh/networkservicemesh/k8s/cmd/nsm-coredns/plugin/fanout"
)

var version string

func init() {
	dnsserver.Directives = append(dnsserver.Directives, "fanout")
}

func main() {
	log.Println("Starting nsm-coredns")
	log.Printf("Version: %v\n", version)
	utils.PrintAllEnv(logrus.StandardLogger())
	defaultConfig := defaultBasicDNSConfig()
	log.Printf("Default dns config: %v", defaultConfig)
	updateResolvConfFile()
	path := parseCorefilePath()
	if env.UseUpdateAPIEnv.GetBooleanOrDefault(false) {
		file := caddyfile.NewCaddyfile(path)
		file.WriteScope(".").Write("log").Write(fmt.Sprintf("fanout %v", strings.Join(defaultConfig.DnsServerIps, " ")))
		err := file.Save()
		fmt.Println(file.String())
		if err != nil {
			log.Println(err.Error())
			os.Exit(2)
		}
		fmt.Println("Starting dns context update server...")
		err = startUpdateServer()
		if err != nil {
			log.Println(err.Error())
			os.Exit(2)
		}
	}
	coremain.Run()
}
