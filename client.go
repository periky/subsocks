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
	cli.Proxys = append(cli.Proxys, cfg.Client.Proxy...)

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
	cliCrt, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}
	serverName, _, _ := net.SplitHostPort(addr)
	config = &tls.Config{
		ServerName:   serverName,
		RootCAs:      x509.NewCertPool(),
		Certificates: []tls.Certificate{cliCrt},
	}
	err = loadCA(caFile, config)
	if err != nil {
		return
	}

	return
}

func loadCA(caFile string, config *tls.Config) error {
	if caFile == "" {
		return errors.New("cafile not provide")
	}
	data, err := ioutil.ReadFile(caFile)
	if err != nil {
		return err
	}
	if !config.RootCAs.AppendCertsFromPEM(data) {
		return errors.New("append certs from pem failed")
	}
	return nil
}
