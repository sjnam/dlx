package util

import (
	"context"
	"runtime"
)

func OrderedProcess[V any, T any](
	ctx context.Context,
	inStream <-chan T,
	doWork func(T) V,
	cnt ...int,
) <-chan V {
	lvl := runtime.NumCPU()
	if len(cnt) > 0 {
		lvl = cnt[0]
	}

	orDone := func(ctx context.Context, c <-chan V) <-chan V {
		ch := make(chan V)
		go func() {
			defer close(ch)
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case ch <- v:
					case <-ctx.Done():
					}
				}
			}
		}()
		return ch
	}

	chanchan := func() <-chan <-chan V {
		chch := make(chan (<-chan V), lvl)
		go func() {
			defer close(chch)
			for v := range inStream {
				ch := make(chan V)
				chch <- ch

				go func() {
					defer close(ch)
					ch <- doWork(v)
				}()
			}
		}()
		return chch
	}

	// bridge-channel
	return func(ctx context.Context, chch <-chan <-chan V) <-chan V {
		vch := make(chan V)
		go func() {
			defer close(vch)
			for {
				var ch <-chan V
				select {
				case maybe, ok := <-chch:
					if !ok {
						return
					}
					ch = maybe
				case <-ctx.Done():
					return
				}
				for v := range orDone(ctx, ch) {
					select {
					case vch <- v:
					case <-ctx.Done():
					}
				}
			}
		}()
		return vch
	}(ctx, chanchan())
}
