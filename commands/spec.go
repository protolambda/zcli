package commands

import (
	"fmt"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/configs"
	"gopkg.in/yaml.v3"
	"os"
)

type SpecOptions struct {
	LegacyConfig        string `ask:"--legacy-config" help:"Eth2 legacy configuration (combined config and presets), name or path to YAML"`
	LegacyConfigChanged bool   `changed:"legacy-config"`

	Config         string `ask:"--config" help:"Eth2 spec configuration, name or path to YAML"`
	Phase0Preset   string `ask:"--phase0-preset" help:"Eth2 phase0 spec preset, name or path to YAML"`
	AltairPreset   string `ask:"--phase0-preset" help:"Eth2 altair spec preset, name or path to YAML"`
	MergePreset    string `ask:"--phase0-preset" help:"Eth2 merge spec preset, name or path to YAML"`
	ShardingPreset string `ask:"--phase0-preset" help:"Eth2 sharding spec preset, name or path to YAML"`

	// TODO: execution engine config for Merge
	// TODO: trusted setup config for Sharding
}

type LegacyConfig struct {
	CONFIG_NAME           string `yaml:"CONFIG_NAME"`
	common.Phase0Preset   `yaml:",inline"`
	common.AltairPreset   `yaml:",inline"`
	common.MergePreset    `yaml:",inline"`
	common.ShardingPreset `yaml:",inline"`
	common.Config         `yaml:",inline"`
}

func (c *SpecOptions) Spec() (*common.Spec, error) {
	var spec common.Spec

	if c.LegacyConfigChanged {
		switch c.LegacyConfig {
		case "mainnet":
			spec = *configs.Mainnet
		case "minimal":
			spec = *configs.Minimal
		default:
			var legacy LegacyConfig
			f, err := os.Open(c.LegacyConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to open legacy config file: %v", err)
			}
			dec := yaml.NewDecoder(f)
			if err := dec.Decode(&legacy); err != nil {
				return nil, fmt.Errorf("failed to decode legacy config: %v", err)
			}
			spec.PRESET_BASE = legacy.CONFIG_NAME
			spec.Phase0Preset = legacy.Phase0Preset
			spec.AltairPreset = legacy.AltairPreset
			spec.MergePreset = legacy.MergePreset
			spec.ShardingPreset = legacy.ShardingPreset
			spec.Config = legacy.Config
		}
	}

	switch c.Config {
	case "mainnet":
		spec.Config = configs.Mainnet.Config
	case "minimal":
		spec.Config = configs.Minimal.Config
	default:
		f, err := os.Open(c.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %v", err)
		}
		dec := yaml.NewDecoder(f)
		if err := dec.Decode(&spec.Config); err != nil {
			return nil, fmt.Errorf("failed to decode config: %v", err)
		}
	}

	switch c.Phase0Preset {
	case "mainnet":
		spec.Phase0Preset = configs.Mainnet.Phase0Preset
	case "minimal":
		spec.Phase0Preset = configs.Minimal.Phase0Preset
	default:
		f, err := os.Open(c.Phase0Preset)
		if err != nil {
			return nil, fmt.Errorf("failed to open phase0 preset file: %v", err)
		}
		dec := yaml.NewDecoder(f)
		if err := dec.Decode(&spec.Phase0Preset); err != nil {
			return nil, fmt.Errorf("failed to decode phase0 preset: %v", err)
		}
	}

	switch c.AltairPreset {
	case "mainnet":
		spec.AltairPreset = configs.Mainnet.AltairPreset
	case "minimal":
		spec.AltairPreset = configs.Minimal.AltairPreset
	default:
		f, err := os.Open(c.AltairPreset)
		if err != nil {
			return nil, fmt.Errorf("failed to open altair preset file: %v", err)
		}
		dec := yaml.NewDecoder(f)
		if err := dec.Decode(&spec.AltairPreset); err != nil {
			return nil, fmt.Errorf("failed to decode altair preset: %v", err)
		}
	}

	switch c.MergePreset {
	case "mainnet":
		spec.MergePreset = configs.Mainnet.MergePreset
	case "minimal":
		spec.MergePreset = configs.Minimal.MergePreset
	default:
		f, err := os.Open(c.MergePreset)
		if err != nil {
			return nil, fmt.Errorf("failed to open merge preset file: %v", err)
		}
		dec := yaml.NewDecoder(f)
		if err := dec.Decode(&spec.MergePreset); err != nil {
			return nil, fmt.Errorf("failed to decode merge preset: %v", err)
		}
	}

	switch c.ShardingPreset {
	case "mainnet":
		spec.ShardingPreset = configs.Mainnet.ShardingPreset
	case "minimal":
		spec.ShardingPreset = configs.Minimal.ShardingPreset
	default:
		f, err := os.Open(c.ShardingPreset)
		if err != nil {
			return nil, fmt.Errorf("failed to open sharding preset file: %v", err)
		}
		dec := yaml.NewDecoder(f)
		if err := dec.Decode(&spec.ShardingPreset); err != nil {
			return nil, fmt.Errorf("failed to decode sharding preset: %v", err)
		}
	}
	return &spec, nil
}

func (c *SpecOptions) Default() {
	c.LegacyConfig = "mainnet"
	c.Config = "mainnet"
	c.Phase0Preset = "mainnet"
	c.AltairPreset = "mainnet"
	c.MergePreset = "mainnet"
	c.ShardingPreset = "mainnet"
}
