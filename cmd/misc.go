package cmd

import (
	"os"
	"os/exec"
)

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

// Run another command, as the last thing the CLI does.
// This is used by commands like system ssh that just execute "ssh" and then exit.
func tailExecProcess(command string, args []string) error {
	// TODO: consider using execve on Unix platforms.

	scpCmd := exec.Command(command, args...)
	scpCmd.Stdin = os.Stdin
	scpCmd.Stdout = os.Stdout
	scpCmd.Stderr = os.Stderr
	runErr := scpCmd.Run()
	if exitErr, ok := runErr.(*exec.ExitError); ok {
		os.Exit(exitErr.ExitCode())
	}

	return runErr
}

type tuple2[T1 any, T2 any] struct {
	Item1 T1
	Item2 T2
}

func makeTuple2[T1 any, T2 any](item1 T1, item2 T2) tuple2[T1, T2] {
	return tuple2[T1, T2]{
		Item1: item1,
		Item2: item2,
	}
}
