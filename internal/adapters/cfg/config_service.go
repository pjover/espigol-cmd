package cfg

import (
	"log"
	"os"
	"time"

	"github.com/pjover/espigol/internal/domain/ports"
	"github.com/spf13/viper"
)

type configService struct{}

func NewConfigService() ports.ConfigService {
	service := &configService{}
	service.Init()
	return service
}

func (c *configService) GetString(key string) string {
	return viper.GetString(key)
}

func (c *configService) SetString(key string, value string) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}

func (c *configService) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func (c *configService) GetTime(key string) time.Time {
	return viper.GetTime(key)
}

func (c *configService) SetTime(key string, value time.Time) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}

func (c *configService) Init() {
	home := c.findHomeDirectory()
	c.searchConfigInRootDirectory()
	c.loadEnvironmentVariables()
	c.readConfigFile()
	c.loadDefaultConfig(home)
}

func (c *configService) findHomeDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cannot find home directory: %s", err)
	}
	return home
}

func (c *configService) searchConfigInRootDirectory() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.SetConfigName("espigol")
	viper.SetConfigType("yaml")
}

func (c *configService) loadEnvironmentVariables() {
	viper.AutomaticEnv()
}

func (c *configService) readConfigFile() {
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("Config file not found, using defaults")
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Fatal error config file: %s", err)
		}
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}

func (c *configService) loadDefaultConfig(home string) {
	for key, value := range defaultValues {
		viper.SetDefault(key, value)
	}
}
