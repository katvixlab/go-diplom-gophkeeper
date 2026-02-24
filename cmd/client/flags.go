package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"strings"
)

const defaultClientConfigPath = "config_c.json"

var (
	connAddr string
	logLevel string
	logFile  string
	confFile string
)

type clientConfig struct {
	ConnAddr string `json:"conn_addr"`
	LogLevel string `json:"log_level"`
	LogFile  string `json:"log_file"`
}

func parseFlags() {
	confFile = resolveConfigPath(defaultClientConfigPath)
	defaults := clientConfig{
		ConnAddr: "localhost:3200",
		LogLevel: "info",
		LogFile:  "logs.log",
	}

	if cfg, err := loadClientConfig(confFile); err == nil {
		if cfg.ConnAddr != "" {
			defaults.ConnAddr = cfg.ConnAddr
		}
		if cfg.LogLevel != "" {
			defaults.LogLevel = cfg.LogLevel
		}
		if cfg.LogFile != "" {
			defaults.LogFile = cfg.LogFile
		}
	}

	flag.StringVar(&confFile, "cfg", confFile, "config file path")
	flag.StringVar(&connAddr, "a", defaults.ConnAddr, "server connection address")
	flag.StringVar(&logLevel, "ll", defaults.LogLevel, "log level")
	flag.StringVar(&logFile, "lf", defaults.LogFile, "log path")
	flag.Parse()

	_ = saveClientConfig(confFile, &clientConfig{
		ConnAddr: connAddr,
		LogLevel: logLevel,
		LogFile:  logFile,
	})
}

func resolveConfigPath(defaultPath string) string {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case strings.HasPrefix(arg, "-cfg="):
			path := strings.TrimSpace(strings.TrimPrefix(arg, "-cfg="))
			if path != "" {
				return path
			}
		case arg == "-cfg" && i+1 < len(args):
			path := strings.TrimSpace(args[i+1])
			if path != "" {
				return path
			}
		}
	}
	return defaultPath
}

func saveClientConfig(path string, cfg *clientConfig) error {
	bytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0o644)
}

func loadClientConfig(path string) (*clientConfig, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg clientConfig
	if err = json.Unmarshal(buf, &cfg); err != nil {
		return nil, err
	}

	if cfg.ConnAddr == "" && cfg.LogLevel == "" && cfg.LogFile == "" {
		return nil, errors.New("empty config")
	}
	return &cfg, nil
}
