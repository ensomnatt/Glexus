package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
  "github.com/pelletier/go-toml/v2"
)

type Config struct {
  VideoDir string 
}

func NewConfig() *Config {
  configDir, err := os.UserConfigDir() 
  if err != nil {
    logrus.Fatalf("failed to get user config file: %v", err)
  }
  appDir := filepath.Join(configDir, "glexus")
  configFile := filepath.Join(appDir, "config.toml")

  _, err = os.Stat(appDir)
  if os.IsNotExist(err) {
    err := os.MkdirAll(appDir, 0755)
    if err != nil {
      logrus.Fatalf("error with creating config dir: %v", err)
    }

    file, err := os.Create(configFile)
    if err != nil {
      log.Fatalf("error with creating config file: %v", err)
    }

    cfg := Config{
      VideoDir: "",
    }

    data, err := toml.Marshal(cfg)
    if err != nil {
      log.Fatalf("error with marshaling default config: %v", err)
    }

    _, err = file.Write(data)
    if err != nil {
      log.Fatalf("error with writing default config to the file: %v", err)
    }
  }

  configData, err := os.ReadFile(configFile)
  cfg := &Config{}
  err = toml.Unmarshal(configData, cfg)
  if err != nil {
    logrus.Fatalf("error with reading config file: %v", err)
  }

  return cfg
}
