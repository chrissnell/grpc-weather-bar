package main

import (
	"io/ioutil"

	"github.com/go-ini/ini"
)

// Config is the base configuraiton object
type Config struct {
	Server ServerConfig
	Format FormatConfig
}

// ServerConfig holds configuration related to the remote gRPC server
type ServerConfig struct {
	Hostname string `ini:"hostname"`
	Port     string `ini:"port"`
}

// FormatConfig holds our output formatting configuration
type FormatConfig struct {
	WxFormat string `ini:"weather-format"`
	WindN    string `ini:"wind-n"`
	WindNE   string `ini:"wind-ne"`
	WindE    string `ini:"wind-e"`
	WindSE   string `ini:"wind-se"`
	WindS    string `ini:"wind-s"`
	WindSW   string `ini:"wind-sw"`
	WindW    string `ini:"wind-w"`
	WindNW   string `ini:"wind-nw"`
}

// NewConfig creates an new config object from the given filename.
func NewConfig(filename string) (*Config, error) {
	c := new(Config)
	cfgFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return &Config{}, err
	}

	cfg, err := ini.Load(cfgFile)
	if err != nil {
		return &Config{}, err
	}

	err = cfg.Section("server").MapTo(&c.Server)
	if err != nil {
		return &Config{}, err
	}
	err = cfg.Section("format").MapTo(&c.Format)
	if err != nil {
		return &Config{}, err
	}

	return c, nil
}
