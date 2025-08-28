package fanin

import "context"

func FanIn[T any](ctx context.Context, chans ...<-chan T) <-chan T {
	out := make(chan T)

	for _, ch := range chans {
		go func(ctx context.Context, c <-chan T) {
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case out <- v:
					}
				}
			}
		}(ctx, ch)
	}
	return out
}
