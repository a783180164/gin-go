package configs

import (
	"bytes"
	_ "embed"
	"gin-go/pkg/file"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

//go:embed dev_config.toml
var devConfig []byte

// Config 定义了你要反序列化的结构
type Config struct {
	BaseConfig struct {
		Port string `toml:"port"`
	}

	Logger struct {
		FilePath          string `toml:"filepath"`
		EnableRequestLog  bool   `toml:"enableRequestLog"`
		EnableResponseLog bool   `toml:"enableResponseLog"`
	}

	JWT struct {
		Secret string `toml:"secret"`
		Name   string `toml:name`
		Hour   int32  `toml:hour`
	}

	MySQL struct {
		Read struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"read"`
		Write struct {
			Addr string `toml:"addr"`
			User string `toml:"user"`
			Pass string `toml:"pass"`
			Name string `toml:"name"`
		} `toml:"write"`
		Base struct {
			MaxOpenConn     int           `toml:"maxOpenConn"`
			MaxIdleConn     int           `toml:"maxIdleConn"`
			ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
		} `toml:"base"`
	} `toml:"mysql"`

	QDSTANT struct {
		Host   string `toml:"host"`
		Port   int    `toml:"port"`
		ApiKey string `toml:"apikey"`
	} `toml:"qdrant"`

	OLLAMA struct {
		Host  string `toml:"host"`
		Port  int    `toml:"port"`
		Model string `toml:"model"`
	} `toml:"model"`
}

var config = new(Config)

func init() {
	// 1. 告诉 Viper 这是 TOML
	viper.SetConfigType("toml")

	// 2. 先读取内嵌的默认配置
	if err := viper.ReadConfig(bytes.NewReader(devConfig)); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	// 3. 准备外部配置文件路径
	cfgDir := "./configs"
	cfgFile := filepath.Join(cfgDir, "configs.toml")

	// 4. 如果文件不存在，就帮你写出一份
	if _, exists := file.IsExists(cfgFile); !exists {
		if err := os.MkdirAll(cfgDir, 0766); err != nil {
			panic(err)
		}
		if err := viper.WriteConfigAs(cfgFile); err != nil {
			panic(err)
		}
	} else {
		// 5. 如果已经有文件，就读它来覆盖内嵌默认值
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	}

	// 6. 开启热加载
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
}

// Get 返回当前最新的配置副本
func Get() Config {
	return *config
}
