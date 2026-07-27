package events

import (
	"context"
	"errors"
	"testing"
)

type warmUpTarget struct {
	called bool
	ctx    context.Context
	err    error
}

func (t *warmUpTarget) WarmUp(ctx context.Context) error {
	t.called = true
	t.ctx = ctx
	return t.err
}

type unrelatedTarget struct{}

var _ WarmUp = (*warmUpTarget)(nil)

func TestTriggerContextDiscoversWarmUpHandler(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "value")
	target := &warmUpTarget{}

	if err := TriggerContext(ctx, unrelatedTarget{}, target); err != nil {
		t.Fatalf("TriggerContext() error = %v", err)
	}
	if !target.called {
		t.Fatal("WarmUp was not called")
	}
	if target.ctx != ctx {
		t.Fatal("WarmUp received a different context")
	}
}

func TestTriggerUsesBackgroundContext(t *testing.T) {
	target := &warmUpTarget{}

	if err := Trigger(target); err != nil {
		t.Fatalf("Trigger() error = %v", err)
	}
	if target.ctx == nil {
		t.Fatal("WarmUp received a nil context")
	}
}

func TestTriggerStopsOnFirstError(t *testing.T) {
	want := errors.New("warm-up failed")
	first := &warmUpTarget{err: want}
	second := &warmUpTarget{}

	if err := Trigger(first, second); !errors.Is(err, want) {
		t.Fatalf("Trigger() error = %v, want %v", err, want)
	}
	if second.called {
		t.Fatal("Trigger continued after the first handler error")
	}
}
