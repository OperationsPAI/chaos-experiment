package guidedcli

import (
	"context"
	"testing"
)

func TestBuildInjectionPodKillHappyPath(t *testing.T) {
	cfg := GuidedConfig{
		System:     "ts",
		SystemType: "ts",
		Namespace:  "ts",
		App:        "ts-order-service",
		ChaosType:  "PodKill",
		Duration:   intPtr(1),
	}

	conf, sysType, err := BuildInjection(context.Background(), cfg)
	if err != nil {
		t.Fatalf("BuildInjection returned error: %v", err)
	}
	if conf.PodKill == nil {
		t.Fatalf("expected PodKill spec to be populated, got %+v", conf)
	}
	_ = sysType
}

func TestBuildInjectionRejectsEmptyConfig(t *testing.T) {
	_, _, err := BuildInjection(context.Background(), GuidedConfig{})
	if err == nil {
		t.Fatal("expected BuildInjection to error on empty config")
	}
}
