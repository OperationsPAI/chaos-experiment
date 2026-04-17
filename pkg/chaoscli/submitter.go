package chaoscli

import (
	"context"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Spec mirrors the friendly fault spec shape accepted by aegisctl inject submit.
type Spec struct {
	Type      string         `yaml:"type" json:"type"`
	Namespace string         `yaml:"namespace" json:"namespace"`
	Target    string         `yaml:"target" json:"target"`
	Duration  string         `yaml:"duration" json:"duration"`
	Params    map[string]any `yaml:"params,omitempty" json:"params,omitempty"`
}

// Submitter handles a fully materialized fault spec.
type Submitter interface {
	Submit(ctx context.Context, spec Spec) error
}

// PrintSubmitter writes specs as YAML, useful for smoke tests and dry runs.
type PrintSubmitter struct {
	Writer io.Writer
}

func NewPrintSubmitter(w io.Writer) *PrintSubmitter {
	if w == nil {
		w = os.Stdout
	}
	return &PrintSubmitter{Writer: w}
}

func (p *PrintSubmitter) Submit(_ context.Context, spec Spec) error {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}
	_, err = p.Writer.Write(data)
	return err
}
