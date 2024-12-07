package util

import (
	"context"
	"runtime"
)

func OrderedProcess[T, V any](
	ctx context.Context,
	inStream <-chan V,
	doWork func(V) T,
	cnt ...int,
) <-chan T {
	lvl := runtime.NumCPU()
	if len(cnt) > 0 {
		lvl = cnt[0]
	}

	orDone := func(ctx context.Context, c <-chan T) <-chan T {
		ch := make(chan T)
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

	chanchan := func() <-chan <-chan T {
		chch := make(chan (<-chan T), lvl)
		go func() {
			defer close(chch)
			for v := range inStream {
				ch := make(chan T)
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
	return func(ctx context.Context, chch <-chan <-chan T) <-chan T {
		vch := make(chan T)
		go func() {
			defer close(vch)
			for {
				var ch <-chan T
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
