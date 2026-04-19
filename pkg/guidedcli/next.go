package guidedcli

import (
	"fmt"
	"strconv"
	"strings"
)

// ApplyNext is an alias for ApplyNextSelection kept for external callers
// (e.g. aegisctl) that prefer the shorter name. Given the current guided
// response and a raw --next value, it returns the updated GuidedConfig.
func ApplyNext(response *GuidedResponse, rawValue string) (GuidedConfig, error) {
	return ApplyNextSelection(response, rawValue)
}

func ApplyNextSelection(response *GuidedResponse, rawValue string) (GuidedConfig, error) {
	if response == nil {
		return GuidedConfig{}, fmt.Errorf("guided response is required")
	}

	field, err := nextSelectionField(response.Next)
	if err != nil {
		return GuidedConfig{}, err
	}

	selection, err := buildNextSelection(field, rawValue)
	if err != nil {
		return GuidedConfig{}, err
	}

	cfg := response.Config
	overlayConfig(&cfg, selection)
	return cfg, nil
}

func nextSelectionField(fields []FieldSpec) (FieldSpec, error) {
	if len(fields) == 0 {
		return FieldSpec{}, fmt.Errorf("the current response has no next-step fields")
	}

	required := make([]FieldSpec, 0, len(fields))
	fallback := make([]FieldSpec, 0, len(fields))
	for _, field := range fields {
		if field.Kind == "group" {
			continue
		}
		fallback = append(fallback, field)
		if field.Required {
			required = append(required, field)
		}
	}

	switch {
	case len(required) == 1:
		return required[0], nil
	case len(required) > 1:
		return FieldSpec{}, fmt.Errorf("the current response needs multiple required selections; use explicit flags instead of --next")
	case len(fallback) == 1:
		return fallback[0], nil
	case len(fallback) > 1:
		return FieldSpec{}, fmt.Errorf("the current response exposes multiple selectable fields; use explicit flags instead of --next")
	default:
		return FieldSpec{}, fmt.Errorf("the current response expects a grouped parameter input; use explicit flags instead of --next")
	}
}

func buildNextSelection(field FieldSpec, rawValue string) (GuidedConfig, error) {
	value := strings.TrimSpace(rawValue)
	if value == "" {
		return GuidedConfig{}, fmt.Errorf("--next requires a non-empty value")
	}

	switch field.Kind {
	case "enum":
		return buildEnumSelection(field, value)
	case "number_range":
		return buildNumberSelection(field, value)
	case "object_ref":
		return buildObjectRefSelection(field, value)
	default:
		return GuidedConfig{}, fmt.Errorf("--next does not support field kind %q", field.Kind)
	}
}

func buildEnumSelection(field FieldSpec, rawValue string) (GuidedConfig, error) {
	value := rawValue
	if len(field.Options) > 0 {
		matched, ok := matchOptionValue(field.Options, rawValue)
		if !ok {
			return GuidedConfig{}, fmt.Errorf("%q is not a valid option for %s", rawValue, field.Name)
		}
		value = matched
	}

	cfg := GuidedConfig{}
	if err := assignStringField(&cfg, field.Name, value); err != nil {
		return GuidedConfig{}, err
	}
	return cfg, nil
}

func buildNumberSelection(field FieldSpec, rawValue string) (GuidedConfig, error) {
	number, err := strconv.Atoi(rawValue)
	if err != nil {
		return GuidedConfig{}, fmt.Errorf("%q is not a valid integer for %s", rawValue, field.Name)
	}
	if field.Min != nil && number < *field.Min {
		return GuidedConfig{}, fmt.Errorf("%s must be at least %d", field.Name, *field.Min)
	}
	if field.Max != nil && number > *field.Max {
		return GuidedConfig{}, fmt.Errorf("%s must be at most %d", field.Name, *field.Max)
	}

	cfg := GuidedConfig{}
	if err := assignNumberField(&cfg, field.Name, number); err != nil {
		return GuidedConfig{}, err
	}
	return cfg, nil
}

func buildObjectRefSelection(field FieldSpec, rawValue string) (GuidedConfig, error) {
	values, ok := matchObjectRef(field, rawValue)
	if !ok {
		var err error
		values, err = parseObjectRef(field, rawValue)
		if err != nil {
			return GuidedConfig{}, err
		}
	}

	cfg := GuidedConfig{}
	for key, value := range values {
		if err := assignStringField(&cfg, key, value); err != nil {
			return GuidedConfig{}, err
		}
	}
	return cfg, nil
}

func matchOptionValue(options []FieldOption, rawValue string) (string, bool) {
	for _, option := range options {
		if rawValue == option.Value || rawValue == option.Label {
			return option.Value, true
		}
	}
	return "", false
}

func matchObjectRef(field FieldSpec, rawValue string) (map[string]string, bool) {
	for _, option := range field.Options {
		if rawValue != option.Value && rawValue != option.Label {
			continue
		}

		values := make(map[string]string, len(field.KeyFields))
		for _, key := range field.KeyFields {
			value, ok := option.Metadata[key].(string)
			if !ok || strings.TrimSpace(value) == "" {
				return nil, false
			}
			values[key] = value
		}
		return values, true
	}
	return nil, false
}

func parseObjectRef(field FieldSpec, rawValue string) (map[string]string, error) {
	switch strings.Join(field.KeyFields, ",") {
	case "http_method,route":
		parts := strings.SplitN(rawValue, " ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected HTTP endpoint in the form '<METHOD> <ROUTE>'")
		}
		return map[string]string{
			"http_method": strings.TrimSpace(parts[0]),
			"route":       strings.TrimSpace(parts[1]),
		}, nil
	case "class,method":
		parts := strings.SplitN(rawValue, "#", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected JVM method in the form '<Class>#<method>'")
		}
		return map[string]string{
			"class":  strings.TrimSpace(parts[0]),
			"method": strings.TrimSpace(parts[1]),
		}, nil
	case "database,table,operation":
		parts := strings.SplitN(rawValue, "/", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("expected database operation in the form '<database>/<table>/<operation>'")
		}
		return map[string]string{
			"database":  strings.TrimSpace(parts[0]),
			"table":     strings.TrimSpace(parts[1]),
			"operation": strings.TrimSpace(parts[2]),
		}, nil
	default:
		return nil, fmt.Errorf("--next does not know how to parse object ref %q", field.Name)
	}
}

func assignStringField(cfg *GuidedConfig, name, value string) error {
	switch name {
	case "system":
		cfg.System = value
	case "system_type":
		cfg.SystemType = value
	case "namespace":
		cfg.Namespace = value
	case "app":
		cfg.App = value
	case "chaos_type":
		cfg.ChaosType = value
	case "container":
		cfg.Container = value
	case "target_service":
		cfg.TargetService = value
	case "domain":
		cfg.Domain = value
	case "class":
		cfg.Class = value
	case "method":
		cfg.Method = value
	case "mutator_config":
		cfg.MutatorConfig = value
	case "route":
		cfg.Route = value
	case "http_method":
		cfg.HTTPMethod = value
	case "database":
		cfg.Database = value
	case "table":
		cfg.Table = value
	case "operation":
		cfg.Operation = value
	case "direction":
		cfg.Direction = value
	case "return_type":
		cfg.ReturnType = value
	case "return_value_opt":
		cfg.ReturnValueOpt = value
	case "exception_opt":
		cfg.ExceptionOpt = value
	case "mem_type":
		cfg.MemType = value
	case "body_type":
		cfg.BodyType = value
	case "replace_method":
		cfg.ReplaceMethod = value
	default:
		return fmt.Errorf("field %q is not a supported string selection", name)
	}
	return nil
}

func assignNumberField(cfg *GuidedConfig, name string, value int) error {
	switch name {
	case "duration":
		cfg.Duration = intPtr(value)
	case "memory_size":
		cfg.MemorySize = intPtr(value)
	case "mem_worker":
		cfg.MemWorker = intPtr(value)
	case "time_offset":
		cfg.TimeOffset = intPtr(value)
	case "cpu_load":
		cfg.CPULoad = intPtr(value)
	case "cpu_worker":
		cfg.CPUWorker = intPtr(value)
	case "latency":
		cfg.Latency = intPtr(value)
	case "correlation":
		cfg.Correlation = intPtr(value)
	case "jitter":
		cfg.Jitter = intPtr(value)
	case "loss":
		cfg.Loss = intPtr(value)
	case "duplicate":
		cfg.Duplicate = intPtr(value)
	case "corrupt":
		cfg.Corrupt = intPtr(value)
	case "rate":
		cfg.Rate = intPtr(value)
	case "limit":
		cfg.Limit = intPtr(value)
	case "buffer":
		cfg.Buffer = intPtr(value)
	case "delay_duration":
		cfg.DelayDuration = intPtr(value)
	case "latency_duration":
		cfg.LatencyDuration = intPtr(value)
	case "latency_ms":
		cfg.LatencyMs = intPtr(value)
	case "cpu_count":
		cfg.CPUCount = intPtr(value)
	case "status_code":
		cfg.StatusCode = intPtr(value)
	default:
		return fmt.Errorf("field %q is not a supported numeric selection", name)
	}
	return nil
}
