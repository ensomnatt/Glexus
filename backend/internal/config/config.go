package config

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/pelletier/go-toml/v2"
	"github.com/sirupsen/logrus"
)

type Config struct {
  configDir string
  appDir string
  configFile string
  VideoFiles []string
  VideoDir string 
}

func NewConfig() *Config {
  configDir, err := os.UserConfigDir() 
  if err != nil {
    logrus.Fatalf("failed to get user config file: %v", err)
  }
  appDir := filepath.Join(configDir, "glexus")
  configFile := filepath.Join(appDir, "config.toml")

  cfg := &Config{
    configDir: configDir,
    appDir: appDir,
    configFile: configFile,
  }

  _, err = os.Stat(appDir)
  if os.IsNotExist(err) {
    cfg.createConfigFile()
  }

  configData, err := os.ReadFile(configFile)
  err = toml.Unmarshal(configData, cfg)
  if err != nil {
    logrus.Fatalf("error with reading config file: %v", err)
  }
  logrus.Info("read the config file")

  cfg.getVideoFiles()

  /*
  cfg { 
    configDir: configDir,
    appDir: appDir,
    configFile: configFile,
    VideoFiles: sorted video files,
    VideoDir: videoDir,
  }
  */

  return cfg
}

func (c *Config) createConfigFile() {
    logrus.Warn("can't find config dir. creating...")
    err := os.MkdirAll(c.appDir, 0755)
    if err != nil {
      logrus.Fatalf("error with creating config dir: %v", err)
    }
    logrus.Info("created config dir")

    file, err := os.Create(c.configFile)
    if err != nil {
      logrus.Fatalf("error with creating config file: %v", err)
    }
    logrus.Info("created config file")

    cfg := Config{
      VideoDir: "",
    }

    data, err := toml.Marshal(cfg)
    if err != nil {
      logrus.Fatalf("error with marshaling default config: %v", err)
    }

    _, err = file.Write(data)
    if err != nil {
      logrus.Fatalf("error with writing default config to the file: %v", err)
    }
    logrus.Info("wrote default config, config file was successfully create! you can edit it")
}

func (c *Config) getVideoFiles() {
  _, err := os.Stat(c.VideoDir)
  if os.IsNotExist(err) {
    logrus.Fatal("video directory didn't found")
  }
  videoFiles, err := os.ReadDir(c.VideoDir)
  if err != nil {
    logrus.Fatalf("failed to get video files: %v", err)
  }

  var videoFilesSorted []string
  for _, file := range videoFiles {
    info, _ := os.Stat(filepath.Join(c.VideoDir, file.Name()))
    if !info.IsDir() {
      videoFilesSorted = append(videoFilesSorted, file.Name())
    }
  }

  re := regexp.MustCompile(`\d+`)

  sort.Slice(videoFilesSorted, func(i, j int) bool {
    num1String := re.FindString(videoFilesSorted[i])
    num2String := re.FindString(videoFilesSorted[j])

    num1, _ := strconv.Atoi(num1String)
    num2, _ := strconv.Atoi(num2String)

    return num1 < num2
  })

  c.VideoFiles = videoFilesSorted
} 
