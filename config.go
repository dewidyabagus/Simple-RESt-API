package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

const defListen = "8001"

var defAllowMethods = []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}

type Config struct {
	Listen       string
	AllowMethods []string
	Whitelist    Whitelist
	DatabaseDSN  string
}

type Whitelist struct {
	Enabled bool
	IpAddr  []string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Load .env file:", err)
	}

	config := &Config{
		Listen:       os.Getenv("LISTEN_ADDR"),
		AllowMethods: strings.Split(os.Getenv("ALLOW_METHODS"), ","),
		DatabaseDSN:  strings.TrimSpace(os.Getenv("DATABASE_DSN")),
	}

	if strings.TrimSpace(config.Listen) == "" {
		config.Listen = defListen
	}

	config.Whitelist.Enabled, _ = strconv.ParseBool(os.Getenv("WHITELIST"))
	if config.Whitelist.Enabled {
		for _, val := range strings.Split(os.Getenv("WHITELIST_IP_ADDR"), ",") {
			config.Whitelist.IpAddr = append(config.Whitelist.IpAddr, strings.TrimSpace(val))
		}
	}

	if len(config.AllowMethods) == 0 || slices.Contains(config.AllowMethods, "*") {
		config.AllowMethods = defAllowMethods
	} else {
		for i := range config.AllowMethods {
			config.AllowMethods[i] = strings.ToUpper(strings.TrimSpace(config.AllowMethods[i]))
		}
	}

	return config
}
