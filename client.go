package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net"

	"github.com/periky/subsocks/client"
	"github.com/periky/subsocks/config"
)

func launchClient(cfg *config.Config) {
	cli := client.NewClient(cfg.Client.Listen)
	cli.Config.Username = cfg.Client.UserName
	cli.Config.Password = cfg.Client.Password
	cli.Config.ServerProtocol = cfg.Client.Protocol
	cli.Config.ServerAddr = cfg.Client.Addr
	cli.Config.HTTPPath = cfg.Http.Path
	cli.Config.WSPath = cfg.Ws.Path
	cli.DefaultProxys = append(cli.DefaultProxys, cfg.Client.Proxy...)
	cli.Proxys = cli.DefaultProxys

	if needsTLS[cfg.Client.Protocol] {
		tlsConfig, err := getClientTLSConfig(cfg.Client.Addr, cfg.Tls.CaFile, cfg.Tls.CertFile, cfg.Tls.KeyFile)
		if err != nil {
			log.Fatalf("Get TLS configuration failed: %s", err)
		}
		cli.TLSConfig = tlsConfig
	}

	if err := cli.Serve(); err != nil {
		log.Fatalf("Launch client failed: %s", err)
	}
}

func getClientTLSConfig(addr, caFile, certFile, keyFile string) (config *tls.Config, err error) {
	rootCAs, err := loadCA(caFile)
	if err != nil {
		return
	}
	cliCrt, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}
	serverName, _, _ := net.SplitHostPort(addr)
	config = &tls.Config{
		ServerName:   serverName,
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{cliCrt},
	}

	return
}

func loadCA(caFile string) (*x509.CertPool, error) {
	if caFile == "" {
		return nil, errors.New("cafile not provide")
	}
	data, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(data) {
		return nil, errors.New("append certs from pem failed")
	}
	return cp, nil
}
