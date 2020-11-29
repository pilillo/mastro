package conf

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Args ... Arguments provided either as env vars or string args
var Args struct {
	ConfigPath string `required:"true" arg:"-c,required"`
}

// Config ... Defines a model for the input config files
type Config struct {
	ConfigType           string               `yaml:"type"`
	Details              Details              `yaml:"details,omitempty"`
	DataSourceDefinition DataSourceDefinition `yaml:"backend"`
}

// Details ... a map where we can place service specific configuration
type Details struct {
	Values map[string]string
}

// UnmarshalYAML is used to unmarshal into map[string]string
// https://www.ribice.ba/golang-yaml-string-map/
func (d *Details) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&d.Values)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func parseCfg(data []byte) (*Config, error) {
	cfg := &Config{}

	err := yaml.Unmarshal(data, &cfg)
	log.Println("Successfully loaded config", cfg.ConfigType)

	return cfg, err
}

func validateCfg(cfg *Config) (*Config, error) {
	// todo add validation of input config
	return cfg, nil
}

// Load ... load configuration from file path
func Load(filename string) *Config {
	if !fileExists(filename) {
		log.Fatalf("Example file does not exist (or is a directory)")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error - %v", err)
	}

	config, err := parseCfg(data)
	if err != nil {
		log.Fatalf("Error - %v", err)
	}
	config, err = validateCfg(config)
	if err != nil {
		log.Fatalf("Error - %v", err)
	}

	return config
}
