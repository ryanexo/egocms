package events

import "context"

type WarmUp interface {
	WarmUp(context.Context) error
}

func Trigger(targets ...any) error {
	return TriggerContext(context.Background(), targets...)
}

func TriggerContext(ctx context.Context, targets ...any) error {
	if ctx == nil {
		ctx = context.Background()
	}

	for _, target := range targets {
		if target == nil {
			continue
		}

		if event, ok := target.(WarmUp); ok {
			if err := event.WarmUp(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}
