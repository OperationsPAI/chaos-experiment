package guidedcli

type GuidedConfig struct {
	System          string `json:"system,omitempty" yaml:"system,omitempty"`
	SystemType      string `json:"system_type,omitempty" yaml:"system_type,omitempty"`
	Namespace       string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	App             string `json:"app,omitempty" yaml:"app,omitempty"`
	ChaosType       string `json:"chaos_type,omitempty" yaml:"chaos_type,omitempty"`
	Container       string `json:"container,omitempty" yaml:"container,omitempty"`
	TargetService   string `json:"target_service,omitempty" yaml:"target_service,omitempty"`
	Domain          string `json:"domain,omitempty" yaml:"domain,omitempty"`
	Class           string `json:"class,omitempty" yaml:"class,omitempty"`
	Method          string `json:"method,omitempty" yaml:"method,omitempty"`
	MutatorConfig   string `json:"mutator_config,omitempty" yaml:"mutator_config,omitempty"`
	Route           string `json:"route,omitempty" yaml:"route,omitempty"`
	HTTPMethod      string `json:"http_method,omitempty" yaml:"http_method,omitempty"`
	Database        string `json:"database,omitempty" yaml:"database,omitempty"`
	Table           string `json:"table,omitempty" yaml:"table,omitempty"`
	Operation       string `json:"operation,omitempty" yaml:"operation,omitempty"`
	Duration        *int   `json:"duration,omitempty" yaml:"duration,omitempty"`
	MemorySize      *int   `json:"memory_size,omitempty" yaml:"memory_size,omitempty"`
	MemWorker       *int   `json:"mem_worker,omitempty" yaml:"mem_worker,omitempty"`
	TimeOffset      *int   `json:"time_offset,omitempty" yaml:"time_offset,omitempty"`
	CPULoad         *int   `json:"cpu_load,omitempty" yaml:"cpu_load,omitempty"`
	CPUWorker       *int   `json:"cpu_worker,omitempty" yaml:"cpu_worker,omitempty"`
	Latency         *int   `json:"latency,omitempty" yaml:"latency,omitempty"`
	Correlation     *int   `json:"correlation,omitempty" yaml:"correlation,omitempty"`
	Jitter          *int   `json:"jitter,omitempty" yaml:"jitter,omitempty"`
	Loss            *int   `json:"loss,omitempty" yaml:"loss,omitempty"`
	Duplicate       *int   `json:"duplicate,omitempty" yaml:"duplicate,omitempty"`
	Corrupt         *int   `json:"corrupt,omitempty" yaml:"corrupt,omitempty"`
	Rate            *int   `json:"rate,omitempty" yaml:"rate,omitempty"`
	Limit           *int   `json:"limit,omitempty" yaml:"limit,omitempty"`
	Buffer          *int   `json:"buffer,omitempty" yaml:"buffer,omitempty"`
	Direction       string `json:"direction,omitempty" yaml:"direction,omitempty"`
	DelayDuration   *int   `json:"delay_duration,omitempty" yaml:"delay_duration,omitempty"`
	LatencyDuration *int   `json:"latency_duration,omitempty" yaml:"latency_duration,omitempty"`
	LatencyMs       *int   `json:"latency_ms,omitempty" yaml:"latency_ms,omitempty"`
	CPUCount        *int   `json:"cpu_count,omitempty" yaml:"cpu_count,omitempty"`
	ReturnType      string `json:"return_type,omitempty" yaml:"return_type,omitempty"`
	ReturnValueOpt  string `json:"return_value_opt,omitempty" yaml:"return_value_opt,omitempty"`
	ExceptionOpt    string `json:"exception_opt,omitempty" yaml:"exception_opt,omitempty"`
	MemType         string `json:"mem_type,omitempty" yaml:"mem_type,omitempty"`
	BodyType        string `json:"body_type,omitempty" yaml:"body_type,omitempty"`
	ReplaceMethod   string `json:"replace_method,omitempty" yaml:"replace_method,omitempty"`
	StatusCode      *int   `json:"status_code,omitempty" yaml:"status_code,omitempty"`
	SaveConfig      bool   `json:"-" yaml:"-"`
	ResetConfig     bool   `json:"-" yaml:"-"`
	Apply           bool   `json:"-" yaml:"-"`
}

type GuidedResponse struct {
	Mode         string                 `json:"mode" yaml:"mode"`
	Stage        string                 `json:"stage" yaml:"stage"`
	Config       GuidedConfig           `json:"config" yaml:"config"`
	Resolved     map[string]any         `json:"resolved,omitempty" yaml:"resolved,omitempty"`
	Next         []FieldSpec            `json:"next,omitempty" yaml:"next,omitempty"`
	Preview      *Preview               `json:"preview,omitempty" yaml:"preview,omitempty"`
	ApplyPayload map[string]any         `json:"apply_payload,omitempty" yaml:"apply_payload,omitempty"`
	Result       map[string]any         `json:"result,omitempty" yaml:"result,omitempty"`
	CanApply     bool                   `json:"can_apply" yaml:"can_apply"`
	Warnings     []string               `json:"warnings,omitempty" yaml:"warnings,omitempty"`
	Errors       []string               `json:"errors,omitempty" yaml:"errors,omitempty"`
	Resources    map[string]any         `json:"resources,omitempty" yaml:"resources,omitempty"`
	Meta         map[string]interface{} `json:"meta,omitempty" yaml:"meta,omitempty"`
}

type Preview struct {
	DisplayConfig   map[string]any `json:"display_config,omitempty" yaml:"display_config,omitempty"`
	Groundtruth     map[string]any `json:"groundtruth,omitempty" yaml:"groundtruth,omitempty"`
	ResourceSummary map[string]any `json:"resource_summary,omitempty" yaml:"resource_summary,omitempty"`
}

type FieldSpec struct {
	Name        string        `json:"name" yaml:"name"`
	Kind        string        `json:"kind" yaml:"kind"`
	Required    bool          `json:"required" yaml:"required"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Options     []FieldOption `json:"options,omitempty" yaml:"options,omitempty"`
	Fields      []FieldSpec   `json:"fields,omitempty" yaml:"fields,omitempty"`
	Min         *int          `json:"min,omitempty" yaml:"min,omitempty"`
	Max         *int          `json:"max,omitempty" yaml:"max,omitempty"`
	Step        *int          `json:"step,omitempty" yaml:"step,omitempty"`
	Default     *int          `json:"default,omitempty" yaml:"default,omitempty"`
	Unit        string        `json:"unit,omitempty" yaml:"unit,omitempty"`
	KeyFields   []string      `json:"key_fields,omitempty" yaml:"key_fields,omitempty"`
}

type FieldOption struct {
	Value       string         `json:"value,omitempty" yaml:"value,omitempty"`
	Label       string         `json:"label,omitempty" yaml:"label,omitempty"`
	Description string         `json:"description,omitempty" yaml:"description,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

type ConfigFile struct {
	Version        int                   `yaml:"version"`
	CurrentContext string                `yaml:"current-context,omitempty"`
	Contexts       map[string]CLIContext `yaml:"contexts,omitempty"`
	GuidedSession  GuidedSession         `yaml:"guided-session,omitempty"`
}

type CLIContext struct {
	Kubeconfig        string `yaml:"kubeconfig,omitempty"`
	KubeContext       string `yaml:"kubecontext,omitempty"`
	Output            string `yaml:"output,omitempty"`
	DefaultSystem     string `yaml:"default-system,omitempty"`
	DefaultSystemType string `yaml:"default-system-type,omitempty"`
	DefaultNamespace  string `yaml:"default-namespace,omitempty"`
}

type GuidedSession struct {
	Config    GuidedConfig `yaml:"config,omitempty"`
	UpdatedAt string       `yaml:"updated_at,omitempty"`
}

func intPtr(v int) *int {
	return &v
}

// NewConfig returns a fresh GuidedConfig pre-populated with the given namespace.
// Intended as a starter constructor for external callers (e.g. aegisctl).
func NewConfig(namespace string) *GuidedConfig {
	return &GuidedConfig{Namespace: namespace}
}
