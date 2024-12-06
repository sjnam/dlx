package util

import (
	"context"
	"runtime"
)

func OrderedProcess[T any](
	ctx context.Context,
	inputStream <-chan string,
	doWork func(string) T,
	cnt ...int, /*optional param*/
) <-chan T {
	orDone := func(
		ctx context.Context,
		c <-chan T,
	) <-chan T {
		valStream := make(chan T)
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
		inputStream <-chan string,
		doWork func(string) T,
		clvl int,
	) <-chan <-chan T {
		chStream := make(chan (<-chan T), clvl)
		go func() {
			defer close(chStream)
			for v := range inputStream {
				stream := make(chan T)
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
		chStream <-chan <-chan T,
	) <-chan T {
		valStream := make(chan T)
		go func() {
			defer close(valStream)
			for {
				var stream <-chan T
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

	clvl := runtime.NumCPU() // concurrency level
	if len(cnt) > 0 {
		clvl = cnt[0]
	}

	return bridge(ctx, chanStream(ctx, inputStream, doWork, clvl))
}
