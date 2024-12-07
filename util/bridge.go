package util

import (
	"context"
	"runtime"
)

func OrderedProcess[T1 any, T2 any](
	ctx context.Context,
	inStream <-chan T2,
	doWork func(T2) T1,
	cnt ...int,
) <-chan T1 {
	lvl := runtime.NumCPU()
	if len(cnt) > 0 {
		lvl = cnt[0]
	}

	orDone := func(ctx context.Context, c <-chan T1) <-chan T1 {
		ch := make(chan T1)
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

	chanchan := func(ctx context.Context) <-chan <-chan T1 {
		chch := make(chan (<-chan T1), lvl)
		go func() {
			defer close(chch)
			for v := range inStream {
				ch := make(chan T1)
				select {
				case <-ctx.Done():
					return
				case chch <- ch:
				}

				go func() {
					defer close(ch)
					select {
					case ch <- doWork(v):
					case <-ctx.Done():
					}
				}()
			}
		}()
		return chch
	}

	// bridge-channel
	return func(ctx context.Context, chch <-chan <-chan T1) <-chan T1 {
		vch := make(chan T1)
		go func() {
			defer close(vch)
			for {
				var ch <-chan T1
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
	}(ctx, chanchan(ctx))
}
