package guidedcli

import "testing"

func TestMergeConfigClearsDownstreamWhenAppChanges(t *testing.T) {
	fileCfg := &ConfigFile{
		GuidedSession: GuidedSession{
			Config: GuidedConfig{
				System:        "ts",
				SystemType:    "ts",
				Namespace:     "ts",
				App:           "ts-auth-service",
				ChaosType:     "JVMRuntimeMutator",
				Class:         "auth.Jwt",
				Method:        "createToken",
				MutatorConfig: "string:reverse",
				Duration:      intPtr(9),
			},
		},
	}

	merged := MergeConfig(fileCfg, GuidedConfig{App: "ts-order-service"})
	if merged.App != "ts-order-service" {
		t.Fatalf("expected app override to apply, got %q", merged.App)
	}
	if merged.ChaosType != "" {
		t.Fatalf("expected chaos type to be cleared, got %q", merged.ChaosType)
	}
	if merged.Class != "" || merged.Method != "" || merged.MutatorConfig != "" {
		t.Fatalf("expected downstream JVM selections to be cleared, got class=%q method=%q mutator=%q", merged.Class, merged.Method, merged.MutatorConfig)
	}
	if merged.Duration == nil || *merged.Duration != 9 {
		t.Fatalf("expected duration to be preserved, got %#v", merged.Duration)
	}
}

func TestMergeConfigResetsRootWhenNamespaceChanges(t *testing.T) {
	fileCfg := &ConfigFile{
		GuidedSession: GuidedSession{
			Config: GuidedConfig{
				System:     "ts",
				SystemType: "ts",
				Namespace:  "ts",
				App:        "ts-auth-service",
				ChaosType:  "PodKill",
				Duration:   intPtr(5),
			},
		},
	}

	merged := MergeConfig(fileCfg, GuidedConfig{Namespace: "ts0"})
	if merged.Namespace != "ts0" {
		t.Fatalf("expected namespace override to apply, got %q", merged.Namespace)
	}
	if merged.System != "" || merged.SystemType != "" {
		t.Fatalf("expected saved root fields to be cleared before normalization, got system=%q systemType=%q", merged.System, merged.SystemType)
	}
	if merged.App != "" || merged.ChaosType != "" {
		t.Fatalf("expected downstream selections to be cleared, got app=%q chaosType=%q", merged.App, merged.ChaosType)
	}
}

func TestApplyNextSelectionUsesRequiredField(t *testing.T) {
	response := &GuidedResponse{
		Config: GuidedConfig{
			System:     "ts",
			SystemType: "ts",
			Namespace:  "ts",
			App:        "ts-auth-service",
		},
		Next: []FieldSpec{{
			Name:     "chaos_type",
			Kind:     "enum",
			Required: true,
			Options: []FieldOption{
				{Value: "PodKill", Label: "PodKill"},
				{Value: "PodFailure", Label: "PodFailure"},
			},
		}},
	}

	cfg, err := ApplyNextSelection(response, "PodKill")
	if err != nil {
		t.Fatalf("ApplyNextSelection returned error: %v", err)
	}
	if cfg.ChaosType != "PodKill" {
		t.Fatalf("expected chaos type to be set, got %q", cfg.ChaosType)
	}
}

func TestApplyNextSelectionParsesObjectRef(t *testing.T) {
	response := &GuidedResponse{
		Config: GuidedConfig{
			System:     "ts",
			SystemType: "ts",
			Namespace:  "ts",
			App:        "ts-auth-service",
			ChaosType:  "HTTPRequestDelay",
		},
		Next: []FieldSpec{{
			Name:      "endpoint",
			Kind:      "object_ref",
			Required:  true,
			KeyFields: []string{"http_method", "route"},
			Options: []FieldOption{{
				Value: "POST /api/v1/orders",
				Label: "POST /api/v1/orders",
				Metadata: map[string]any{
					"http_method": "POST",
					"route":       "/api/v1/orders",
				},
			}},
		}},
	}

	cfg, err := ApplyNextSelection(response, "POST /api/v1/orders")
	if err != nil {
		t.Fatalf("ApplyNextSelection returned error: %v", err)
	}
	if cfg.HTTPMethod != "POST" || cfg.Route != "/api/v1/orders" {
		t.Fatalf("expected endpoint selection to be populated, got method=%q route=%q", cfg.HTTPMethod, cfg.Route)
	}
}

func TestApplyNextSelectionRejectsGroupedStage(t *testing.T) {
	response := &GuidedResponse{
		Config: GuidedConfig{ChaosType: "CPUStress"},
		Next: []FieldSpec{{
			Name:     "params",
			Kind:     "group",
			Required: true,
			Fields: []FieldSpec{
				requiredNumberField("cpu_load", "CPU load", 1, 100, 1, "%"),
			},
		}},
	}

	if _, err := ApplyNextSelection(response, "80"); err == nil {
		t.Fatal("expected grouped stage to reject --next")
	}
}
