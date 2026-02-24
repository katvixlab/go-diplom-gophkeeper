package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"strings"
)

const defaultServerConfigPath = "config_s.json"

var (
	srvAddr  string
	logLevel string
	dbFile   string
	crtFile  string
	confFile string
)

type serverConfig struct {
	SrvAddr       string `json:"conn_addr"`
	LogLevel      string `json:"log_level"`
	CrtFile       string `json:"crt_file"`
	DBFile        string `json:"db_file"`
	LegacyDBField string `json:"log_file,omitempty"`
}

func parseFlags() {
	confFile = resolveConfigPath(defaultServerConfigPath)
	defaults := serverConfig{
		SrvAddr:  "localhost:3200",
		LogLevel: "info",
		DBFile:   "test.db",
		CrtFile:  "private.pem",
	}

	if cfg, err := loadServerConfig(confFile); err == nil {
		if cfg.SrvAddr != "" {
			defaults.SrvAddr = cfg.SrvAddr
		}
		if cfg.LogLevel != "" {
			defaults.LogLevel = cfg.LogLevel
		}
		if cfg.CrtFile != "" {
			defaults.CrtFile = cfg.CrtFile
		}
		if cfg.DBFile != "" {
			defaults.DBFile = cfg.DBFile
		}
	}

	flag.StringVar(&confFile, "cfg", confFile, "config file path")
	flag.StringVar(&srvAddr, "a", defaults.SrvAddr, "server address")
	flag.StringVar(&logLevel, "ll", defaults.LogLevel, "log level")
	flag.StringVar(&dbFile, "db", defaults.DBFile, "db path")
	flag.StringVar(&crtFile, "crt", defaults.CrtFile, "certificate x509 path")
	flag.Parse()

	_ = saveServerConfig(confFile, &serverConfig{
		SrvAddr:  srvAddr,
		LogLevel: logLevel,
		DBFile:   dbFile,
		CrtFile:  crtFile,
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

func saveServerConfig(path string, cfg *serverConfig) error {
	bytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0o644)
}

func loadServerConfig(path string) (*serverConfig, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg serverConfig
	if err = json.Unmarshal(buf, &cfg); err != nil {
		return nil, err
	}

	if cfg.DBFile == "" && cfg.LegacyDBField != "" {
		cfg.DBFile = cfg.LegacyDBField
	}

	if cfg.SrvAddr == "" && cfg.LogLevel == "" && cfg.CrtFile == "" && cfg.DBFile == "" {
		return nil, errors.New("empty config")
	}

	return &cfg, nil
}
