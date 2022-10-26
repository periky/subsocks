package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Client *Client `toml:"client"`
	Server *Server `toml:"server"`
	Http   Http    `toml:"http"`
	Ws     Ws      `toml:"ws"`
	Tls    Tls     `toml:"tls"`
}

type Client struct {
	Listen   string   `toml:"listen"`
	Protocol string   `toml:"protocol"`
	UserName string   `toml:"username"`
	Password string   `toml:"password"`
	Addr     string   `toml:"address"`
	Proxy    []string `toml:"proxy"`
}

type Server struct {
	Protocol string `toml:"protocol"`
	Addr     string `toml:"listen"`
}

type Http struct {
	Path string `toml:"path"`
}

type Ws struct {
	Path string `toml:"path"`
	// 仅支持服务端配置生效
	Compress bool `toml:"compress"`
}

type Tls struct {
	CaFile   string `toml:"ca"`
	CertFile string `toml:"cert"`
	KeyFile  string `toml:"key"`
}

func MustParse(filepath string) *Config {
	content, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	cfg := &Config{}
	if err := toml.Unmarshal(content, cfg); err != nil {
		panic(err)
	}

	return cfg
}
