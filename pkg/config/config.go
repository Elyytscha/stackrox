package config

import (
	"os"

	"github.com/ghodss/yaml"
	"github.com/stackrox/rox/pkg/errorhelpers"
)

const (
	path = "/etc/stackrox/central-config.yaml"
)

var (
	defaultBucketFillFraction = 0.5
	defaultCompactionState    = true
)

// Compaction defines the compaction configuration
type Compaction struct {
	Enabled               *bool    `yaml:"enabled"`
	BucketFillFraction    *float64 `yaml:"bucketFillFraction"`
	FreeFractionThreshold *float64 `yaml:"freeFractionThreshold"`
}

func (c *Compaction) applyDefaults() {
	if c.BucketFillFraction == nil {
		c.BucketFillFraction = &defaultBucketFillFraction
	}
	if c.Enabled == nil {
		c.Enabled = &defaultCompactionState
	}
}

// validate must be called after apply defaults
func (c *Compaction) validate() error {
	errorList := errorhelpers.NewErrorList("validating compaction")
	if *c.BucketFillFraction <= 0 || *c.BucketFillFraction > 1.0 {
		errorList.AddString("fill fraction must be greater than 0 and less than or equal to 1")
	}
	if c.FreeFractionThreshold != nil && (*c.FreeFractionThreshold <= 0 || *c.FreeFractionThreshold > 1.0) {
		errorList.AddString("compaction threshold fraction must be greater than 0 and less than or equal to 1")
	}
	return errorList.ToError()
}

// Maintenance defines the maintenance functions to use when Central starts
type Maintenance struct {
	SafeMode   bool       `yaml:"safeMode"`
	Compaction Compaction `yaml:"compaction"`
	// Allowing force rollback only to this version.
	ForceRollbackVersion string `yaml:"forceRollbackVersion"`
}

func (m *Maintenance) applyDefaults() {
	m.Compaction.applyDefaults()
}

func (m *Maintenance) validate() error {
	errorList := errorhelpers.NewErrorList("validating maintenance")
	if err := m.Compaction.validate(); err != nil {
		errorList.AddError(err)
	}
	return errorList.ToError()
}

// Config defines the configuration for Central
type Config struct {
	Maintenance Maintenance `yaml:"maintenance"`
}

func (c *Config) applyDefaults() {
	c.Maintenance.applyDefaults()
}

func (c *Config) validate() error {
	errorList := errorhelpers.NewErrorList("validating config")
	if err := c.Maintenance.validate(); err != nil {
		errorList.AddError(err)
	}
	return errorList.ToError()
}

// ReadConfig reads the configuration file
func ReadConfig() (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			conf := new(Config)
			conf.applyDefaults()
			return conf, nil
		}
		return nil, err
	}
	var conf Config
	if err := yaml.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}
	conf.applyDefaults()
	if err := conf.validate(); err != nil {
		return nil, err
	}
	return &conf, nil
}
