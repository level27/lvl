package cmd

// Finds the index of an element within a slice. Returns -1 if the element is not present.
func indexOf[E comparable](slice []E, value E) int {
	// Replace with golang.org/x/exp/slices/Index if/when that ever becomes non-experimental.
	for idx, v := range slice {
		if v == value {
			return idx
		}
	}

	return -1
}

type resultPair[T any] struct {
	Result T
	Error  error
}

func resultMake[T any](result T, err error) resultPair[T] {
	return resultPair[T]{Result: result, Error: err}
}

func taskRun[T any](task func() (T, error)) <-chan resultPair[T] {
	channel := make(chan resultPair[T])

	go func() {
		channel <- resultMake(task())
		close(channel)
	}()

	return channel
}

func taskRunVoid(task func() error) <-chan error {
	channel := make(chan error)

	go func() {
		channel <- task()
		close(channel)
	}()

	return channel
}
