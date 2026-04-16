package guidedcli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determine home directory: %w", err)
	}
	return filepath.Join(home, ".chaos-exp", "config.yaml"), nil
}

func LoadConfig(path string) (*ConfigFile, error) {
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ConfigFile{
				Version:        1,
				CurrentContext: "default",
				Contexts:       map[string]CLIContext{"default": {}},
			}, nil
		}
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg ConfigFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}
	if cfg.Version == 0 {
		cfg.Version = 1
	}
	if cfg.Contexts == nil {
		cfg.Contexts = map[string]CLIContext{"default": {}}
	}
	if cfg.CurrentContext == "" {
		cfg.CurrentContext = "default"
	}
	return &cfg, nil
}

func SaveConfig(path string, cfg *ConfigFile, snapshot GuidedConfig) error {
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return err
		}
	}

	if cfg == nil {
		cfg = &ConfigFile{Version: 1}
	}
	if cfg.Version == 0 {
		cfg.Version = 1
	}
	cfg.GuidedSession = GuidedSession{
		Config:    snapshot,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config file: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}
	return nil
}

func MergeConfig(fileCfg *ConfigFile, cliCfg GuidedConfig) GuidedConfig {
	merged := GuidedConfig{}
	if fileCfg != nil {
		merged = fileCfg.GuidedSession.Config
		if ctx, ok := fileCfg.Contexts[fileCfg.CurrentContext]; ok {
			if merged.System == "" {
				merged.System = ctx.DefaultSystem
			}
			if merged.SystemType == "" {
				merged.SystemType = ctx.DefaultSystemType
			}
			if merged.Namespace == "" {
				merged.Namespace = ctx.DefaultNamespace
			}
		}
	}

	overlayConfig(&merged, cliCfg)
	return merged
}

func overlayConfig(dst *GuidedConfig, src GuidedConfig) {
	if src.System != "" {
		dst.System = src.System
	}
	if src.SystemType != "" {
		dst.SystemType = src.SystemType
	}
	if src.Namespace != "" {
		dst.Namespace = src.Namespace
	}
	if src.App != "" {
		dst.App = src.App
	}
	if src.ChaosType != "" {
		dst.ChaosType = src.ChaosType
	}
	if src.Container != "" {
		dst.Container = src.Container
	}
	if src.TargetService != "" {
		dst.TargetService = src.TargetService
	}
	if src.Domain != "" {
		dst.Domain = src.Domain
	}
	if src.Class != "" {
		dst.Class = src.Class
	}
	if src.Method != "" {
		dst.Method = src.Method
	}
	if src.MutatorConfig != "" {
		dst.MutatorConfig = src.MutatorConfig
	}
	if src.Route != "" {
		dst.Route = src.Route
	}
	if src.HTTPMethod != "" {
		dst.HTTPMethod = src.HTTPMethod
	}
	if src.Database != "" {
		dst.Database = src.Database
	}
	if src.Table != "" {
		dst.Table = src.Table
	}
	if src.Operation != "" {
		dst.Operation = src.Operation
	}
	if src.Duration != nil {
		dst.Duration = src.Duration
	}
	if src.CPULoad != nil {
		dst.CPULoad = src.CPULoad
	}
	if src.CPUWorker != nil {
		dst.CPUWorker = src.CPUWorker
	}
	if src.Latency != nil {
		dst.Latency = src.Latency
	}
	if src.Correlation != nil {
		dst.Correlation = src.Correlation
	}
	if src.Jitter != nil {
		dst.Jitter = src.Jitter
	}
	if src.Direction != "" {
		dst.Direction = src.Direction
	}
	if src.DelayDuration != nil {
		dst.DelayDuration = src.DelayDuration
	}
	if src.LatencyDuration != nil {
		dst.LatencyDuration = src.LatencyDuration
	}
	dst.SaveConfig = src.SaveConfig
	dst.ResetConfig = src.ResetConfig
	dst.Apply = src.Apply
}
