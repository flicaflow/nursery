package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/flicaflow/nursery"
)

func main() {

	n := nursery.New(context.Background())

	// create channels
	workChan := make(chan int)
	resultChan := make(chan string)

	// Schedule tasks connected by channels
	n.Go(newProducer(workChan))
	n.Go(newWorker(workChan, resultChan))
	n.Go(newCollector(resultChan))

	// wait for the pipeline to complete
	err := n.Wait()
	if err != nil {
		fmt.Printf("Got error: %v", err)
		n.Cancel()
		return
	}

	fmt.Println("Done!")
}

func newProducer(workChan chan<- int) nursery.Task {
	return func(ctx context.Context) error {
		defer close(workChan)
		for i := 0; i < 10; i++ {
			// handle context cancle events to finish the pipeline early
			select {
			case <-ctx.Done():
				fmt.Println("Producer: context canceled")
				return errors.New("Producer: context canceled")
			case workChan <- i:
			}
		}
		return nil
	}
}

func newWorker(workChan <-chan int, resultChan chan<- string) nursery.Task {
	return func(ctx context.Context) error {
		defer close(resultChan)
		for work := range workChan {
			resultChan <- strconv.Itoa(work * work)
		}
		return nil
	}
}

func newCollector(resultChan <-chan string) nursery.Task {
	return func(ctx context.Context) error {
		for result := range resultChan {
			fmt.Printf("Collected Result: %s\n", result)
		}
		return nil
	}
}
