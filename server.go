package main

import (
	"crypto/tls"
	"log"

	"github.com/periky/subsocks/config"
	"github.com/periky/subsocks/server"
)

func launchServer(cfg *config.Config) {
	ser := server.NewServer(cfg.Server.Protocol, cfg.Server.Addr)
	ser.Config.HTTPPath = cfg.Http.Path
	ser.Config.WSPath = cfg.Ws.Path
	ser.Config.WSCompress = cfg.Ws.Compress

	if needsTLS[cfg.Server.Protocol] {
		tlsConfig, err := getServerTLSConfig(cfg.Tls.CaFile, cfg.Tls.CertFile, cfg.Tls.KeyFile)
		if err != nil {
			log.Fatalf("Get TLS configuration failed: %s", err)
		}
		ser.TLSConfig = tlsConfig
	}

	if err := ser.Serve(); err != nil {
		log.Fatalf("Launch server failed: %s", err)
	}
}

func getServerTLSConfig(caFile, certFile, keyFile string) (config *tls.Config, err error) {
	rootCAs, err := loadCA(caFile)
	if err != nil {
		return
	}

	cliCrt, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}

	config = &tls.Config{
		RootCAs:            rootCAs,
		ClientCAs:          rootCAs,
		Certificates:       []tls.Certificate{cliCrt},
		InsecureSkipVerify: false,
		ClientAuth:         tls.RequireAndVerifyClientCert,
	}

	return
}
