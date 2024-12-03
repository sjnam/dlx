package util

import (
	"context"
)

func OrderedProcess(
	ctx context.Context,
	valStream <-chan interface{},
	doWork func(interface{}) interface{},
	cnt int,
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

	resultStream := func(
		ctx context.Context,
		valStream <-chan interface{},
		doWork func(interface{}) interface{},
		cnt int,
	) <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}), cnt)
		go func() {
			defer close(chanStream)
			for v := range valStream {
				ch := make(chan interface{})
				select {
				case <-ctx.Done():
					return
				case chanStream <- ch:
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
		return chanStream
	}

	bridge := func(
		ctx context.Context,
		chanStream <-chan <-chan interface{},
	) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				var stream <-chan interface{}
				select {
				case maybeStream, ok := <-chanStream:
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

	return bridge(ctx, resultStream(ctx, valStream, doWork, cnt))
}
