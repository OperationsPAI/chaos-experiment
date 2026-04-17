package chaoscli

import (
	"bytes"
	"context"
	"testing"

	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestDirectSubmitterBuildsNetworkDelayCRD(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := chaosmeshv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatalf("AddToScheme(chaosmesh): %v", err)
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("AddToScheme(corev1): %v", err)
	}

	k8sClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	submitter := NewDirectSubmitter(k8sClient, &bytes.Buffer{})

	spec := Spec{
		Type:      "NetworkDelay",
		Namespace: "ts",
		Target:    "frontend",
		Duration:  "1m",
		Params: map[string]any{
			"target_service": "checkout",
			"latency":        120,
			"correlation":    50,
			"jitter":         10,
			"direction":      "both",
		},
	}

	if err := submitter.Submit(context.Background(), spec); err != nil {
		t.Fatalf("Submit() error = %v", err)
	}

	var list chaosmeshv1alpha1.NetworkChaosList
	if err := k8sClient.List(context.Background(), &list, client.InNamespace("ts")); err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("expected 1 NetworkChaos, got %d", len(list.Items))
	}

	item := list.Items[0]
	if item.Spec.Action != chaosmeshv1alpha1.DelayAction {
		t.Fatalf("expected action %q, got %q", chaosmeshv1alpha1.DelayAction, item.Spec.Action)
	}
	if got := item.Spec.Selector.LabelSelectors["app"]; got != "frontend" {
		t.Fatalf("expected source app frontend, got %q", got)
	}
	if item.Spec.Target == nil {
		t.Fatal("expected target selector")
	}
	if got := item.Spec.Target.Selector.LabelSelectors["app"]; got != "checkout" {
		t.Fatalf("expected target app checkout, got %q", got)
	}
	if item.Spec.Direction != chaosmeshv1alpha1.Both {
		t.Fatalf("expected direction both, got %q", item.Spec.Direction)
	}
}
