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
	if hasRootOverride(src) {
		clearRootSelection(dst)
	}
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
		if dst.App != src.App {
			clearFromApp(dst)
		}
		dst.App = src.App
	}
	if src.ChaosType != "" {
		if dst.ChaosType != src.ChaosType {
			clearFromChaosType(dst)
		}
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

	if hasMethodOverride(src) {
		if methodSelectionChanged(*dst, src) {
			clearMethodSelection(dst)
		}
		if src.Class != "" {
			dst.Class = src.Class
		}
		if src.Method != "" {
			dst.Method = src.Method
		}
	}
	if src.MutatorConfig != "" {
		dst.MutatorConfig = src.MutatorConfig
	}

	if hasEndpointOverride(src) {
		if endpointSelectionChanged(*dst, src) {
			clearEndpointSelection(dst)
		}
		if src.Route != "" {
			dst.Route = src.Route
		}
		if src.HTTPMethod != "" {
			dst.HTTPMethod = src.HTTPMethod
		}
	}

	if hasDatabaseOverride(src) {
		if databaseSelectionChanged(*dst, src) {
			clearDatabaseSelection(dst)
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
	}
	if src.Duration != nil {
		dst.Duration = src.Duration
	}
	if src.MemorySize != nil {
		dst.MemorySize = src.MemorySize
	}
	if src.MemWorker != nil {
		dst.MemWorker = src.MemWorker
	}
	if src.TimeOffset != nil {
		dst.TimeOffset = src.TimeOffset
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
	if src.Loss != nil {
		dst.Loss = src.Loss
	}
	if src.Duplicate != nil {
		dst.Duplicate = src.Duplicate
	}
	if src.Corrupt != nil {
		dst.Corrupt = src.Corrupt
	}
	if src.Rate != nil {
		dst.Rate = src.Rate
	}
	if src.Limit != nil {
		dst.Limit = src.Limit
	}
	if src.Buffer != nil {
		dst.Buffer = src.Buffer
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
	if src.LatencyMs != nil {
		dst.LatencyMs = src.LatencyMs
	}
	if src.CPUCount != nil {
		dst.CPUCount = src.CPUCount
	}
	if src.ReturnType != "" {
		dst.ReturnType = src.ReturnType
	}
	if src.ReturnValueOpt != "" {
		dst.ReturnValueOpt = src.ReturnValueOpt
	}
	if src.ExceptionOpt != "" {
		dst.ExceptionOpt = src.ExceptionOpt
	}
	if src.MemType != "" {
		dst.MemType = src.MemType
	}
	if src.BodyType != "" {
		dst.BodyType = src.BodyType
	}
	if src.ReplaceMethod != "" {
		dst.ReplaceMethod = src.ReplaceMethod
	}
	if src.StatusCode != nil {
		dst.StatusCode = src.StatusCode
	}
	dst.SaveConfig = src.SaveConfig
	dst.ResetConfig = src.ResetConfig
	dst.Apply = src.Apply
}

func hasRootOverride(cfg GuidedConfig) bool {
	return cfg.System != "" || cfg.SystemType != "" || cfg.Namespace != ""
}

func hasMethodOverride(cfg GuidedConfig) bool {
	return cfg.Class != "" || cfg.Method != ""
}

func hasEndpointOverride(cfg GuidedConfig) bool {
	return cfg.Route != "" || cfg.HTTPMethod != ""
}

func hasDatabaseOverride(cfg GuidedConfig) bool {
	return cfg.Database != "" || cfg.Table != "" || cfg.Operation != ""
}

func clearRootSelection(cfg *GuidedConfig) {
	duration := cfg.Duration
	*cfg = GuidedConfig{}
	cfg.Duration = duration
}

func clearFromApp(cfg *GuidedConfig) {
	cfg.App = ""
	clearFromChaosType(cfg)
}

func clearFromChaosType(cfg *GuidedConfig) {
	cfg.ChaosType = ""
	clearTypeSelections(cfg)
	clearTypeParameters(cfg)
}

func clearTypeSelections(cfg *GuidedConfig) {
	cfg.Container = ""
	cfg.TargetService = ""
	cfg.Domain = ""
	cfg.Class = ""
	cfg.Method = ""
	cfg.MutatorConfig = ""
	cfg.Route = ""
	cfg.HTTPMethod = ""
	cfg.Database = ""
	cfg.Table = ""
	cfg.Operation = ""
}

func clearTypeParameters(cfg *GuidedConfig) {
	cfg.MemorySize = nil
	cfg.MemWorker = nil
	cfg.TimeOffset = nil
	cfg.CPULoad = nil
	cfg.CPUWorker = nil
	cfg.Latency = nil
	cfg.Correlation = nil
	cfg.Jitter = nil
	cfg.Loss = nil
	cfg.Duplicate = nil
	cfg.Corrupt = nil
	cfg.Rate = nil
	cfg.Limit = nil
	cfg.Buffer = nil
	cfg.Direction = ""
	cfg.DelayDuration = nil
	cfg.LatencyDuration = nil
	cfg.LatencyMs = nil
	cfg.CPUCount = nil
	cfg.ReturnType = ""
	cfg.ReturnValueOpt = ""
	cfg.ExceptionOpt = ""
	cfg.MemType = ""
	cfg.BodyType = ""
	cfg.ReplaceMethod = ""
	cfg.StatusCode = nil
}

func clearMethodSelection(cfg *GuidedConfig) {
	cfg.Class = ""
	cfg.Method = ""
	cfg.MutatorConfig = ""
}

func clearEndpointSelection(cfg *GuidedConfig) {
	cfg.Route = ""
	cfg.HTTPMethod = ""
	cfg.ReplaceMethod = ""
}

func clearDatabaseSelection(cfg *GuidedConfig) {
	cfg.Database = ""
	cfg.Table = ""
	cfg.Operation = ""
}

func methodSelectionChanged(current GuidedConfig, incoming GuidedConfig) bool {
	if incoming.Class != "" && incoming.Class != current.Class {
		return true
	}
	if incoming.Method != "" && incoming.Method != current.Method {
		return true
	}
	return false
}

func endpointSelectionChanged(current GuidedConfig, incoming GuidedConfig) bool {
	if incoming.Route != "" && incoming.Route != current.Route {
		return true
	}
	if incoming.HTTPMethod != "" && incoming.HTTPMethod != current.HTTPMethod {
		return true
	}
	return false
}

func databaseSelectionChanged(current GuidedConfig, incoming GuidedConfig) bool {
	if incoming.Database != "" && incoming.Database != current.Database {
		return true
	}
	if incoming.Table != "" && incoming.Table != current.Table {
		return true
	}
	if incoming.Operation != "" && incoming.Operation != current.Operation {
		return true
	}
	return false
}
