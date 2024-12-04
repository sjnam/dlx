package util

import (
	"context"
)

func OrderedProcess(
	ctx context.Context,
	inputStream <-chan interface{},
	doWork func(interface{}) interface{},
	cnt ...int,
) <-chan interface{} {
	orDone := func(
		ctx context.Context,
		c <-chan interface{},
	) <-chan interface{} {
		valStream := make(chan interface{})
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
		inputStream <-chan interface{},
		doWork func(interface{}) interface{},
		clvl int,
	) <-chan <-chan interface{} {
		chStream := make(chan (<-chan interface{}), clvl)
		go func() {
			defer close(chStream)
			for v := range inputStream {
				ch := make(chan interface{})
				select {
				case <-ctx.Done():
					return
				case chStream <- ch:
				}

				go func(v interface{}) {
					defer close(ch)
					select {
					case <-ctx.Done():
						return
					default:
						ch <- doWork(v)
					}
				}(v)
			}
		}()
		return chStream
	}

	bridge := func(
		ctx context.Context,
		chStream <-chan <-chan interface{},
	) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				var stream <-chan interface{}
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

	clvl := 8 // concurrency level
	if len(cnt) > 0 {
		clvl = cnt[0]
	}
	return bridge(ctx, chanStream(ctx, inputStream, doWork, clvl))
}
