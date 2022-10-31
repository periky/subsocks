package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/periky/subsocks/config"
	"github.com/periky/subsocks/utils"
)

func main() {
	var configPath string
	var showVersion bool
	flag.StringVar(&configPath, "c", "./config.toml", "configuration file, default to 'config.toml'")
	flag.BoolVar(&showVersion, "v", false, "show version information")
	flag.Parse()

	if showVersion {
		fmt.Println("Subsocks", Version)
		return
	}

	config := config.MustParse(configPath)
	log.Printf("Load configuration complete: %s", configPath)

	// pprof
	utils.PProf()

	if config.Client != nil {
		launchClient(config)
	}

	if config.Server != nil {
		launchServer(config)
	}
}
