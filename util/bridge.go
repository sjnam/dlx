package util

import (
	"context"
	"runtime"
)

func OrderedProcess[T1 any, T2 any](
	ctx context.Context,
	inputStream <-chan T2,
	doWork func(T2) T1,
	cnt ...int, /*optional param*/
) <-chan T1 {
	clvl := runtime.NumCPU() // concurrency level
	if len(cnt) > 0 {
		clvl = cnt[0]
	}

	orDone := func(
		ctx context.Context,
		c <-chan T1,
	) <-chan T1 {
		valStream := make(chan T1)
		go func() {
			defer close(valStream)
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case valStream <- v:
					case <-ctx.Done():
					}
				}
			}
		}()
		return valStream
	}

	chanStream := func(
		ctx context.Context,
	) <-chan <-chan T1 {
		chStream := make(chan (<-chan T1), clvl)
		go func() {
			defer close(chStream)
			for v := range inputStream {
				stream := make(chan T1)
				select {
				case <-ctx.Done():
					return
				case chStream <- stream:
				}

				go func() {
					defer close(stream)
					select {
					case stream <- doWork(v):
					case <-ctx.Done():
					}
				}()
			}
		}()
		return chStream
	}

	bridge := func(
		ctx context.Context,
		chStream <-chan <-chan T1,
	) <-chan T1 {
		valStream := make(chan T1)
		go func() {
			defer close(valStream)
			for {
				var stream <-chan T1
				select {
				case maybeStream, ok := <-chStream:
					if !ok {
						return
					}
					stream = maybeStream
				case <-ctx.Done():
					return
				}
				for val := range orDone(ctx, stream) {
					select {
					case valStream <- val:
					case <-ctx.Done():
					}
				}
			}
		}()
		return valStream
	}

	return bridge(ctx, chanStream(ctx))
}
