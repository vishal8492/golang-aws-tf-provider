package library

import (
	"context"
	"github.com/hashicorp/terraform-exec/tfexec"
)

// Terraform : terraform wrapper for basic ops
type Terraform interface {
	Init(ctx context.Context) error
	Apply(ctx context.Context) error
	Destroy(ctx context.Context) error
}

type terraform struct {
	provider *tfexec.Terraform
}

func (t *terraform) Init(ctx context.Context) error {
	return t.provider.Init(ctx)
}

func (t *terraform) Apply(ctx context.Context) error {
	return t.provider.Apply(ctx)
}
func (t *terraform) Destroy(ctx context.Context) error {
	return t.provider.Destroy(ctx)
}

func New(provider *tfexec.Terraform) Terraform {
	return &terraform{provider: provider}
}
