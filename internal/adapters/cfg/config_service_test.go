package cfg

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestConfigService(t *testing.T) {
	// Setup a temporary configuration file
	tmpDir := t.TempDir()
	viper.Reset()

	configPath := path.Join(tmpDir, "config_test.yaml")
	viper.SetConfigFile(configPath)

	t.Run("Test Defaults", func(t *testing.T) {
		viper.Reset()

		svc := NewConfigService().(*configService)

		// Verify defaults set in Init()
		val := svc.GetString("db.server")
		if val != "mongodb://localhost:27017" {
			t.Errorf("Expected default db.server to be 'mongodb://localhost:27017', got '%s'", val)
		}
	})

	t.Run("Test SetString and GetString", func(t *testing.T) {
		viper.Reset()

		// Initialize service first (which will fail to find config file, using defaults)
		svc := NewConfigService().(*configService)

		// Set the explicit config file for the test
		viper.SetConfigFile(configPath)

		key := "test.key"
		val := "test_value"

		// Ensure file exists for WriteConfig
		f, _ := os.Create(configPath)
		f.Close()

		err := svc.SetString(key, val)
		if err != nil {
			t.Fatalf("SetString failed: %v", err)
		}

		got := svc.GetString(key)
		if got != val {
			t.Errorf("Expected '%s', got '%s'", val, got)
		}
	})

	t.Run("Test SetTime and GetTime", func(t *testing.T) {
		viper.Reset()

		svc := NewConfigService()

		viper.SetConfigFile(configPath)
		f, _ := os.Create(configPath)
		f.Close()

		key := "test.time"
		now := time.Now().Truncate(time.Second)

		err := svc.SetTime(key, now)
		if err != nil {
			t.Fatalf("SetTime failed: %v", err)
		}

		got := svc.GetTime(key)

		// Force UTC for comparison if needed, but viper usually handles round trip okay enough for this basic test
		if !got.Equal(now) {
			t.Errorf("Expected '%v', got '%v'", now, got)
		}
	})
}
