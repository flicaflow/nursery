package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/flicaflow/nursery"
)

func main() {

	// Create a new nursery
	n := nursery.New(context.Background())

	// schedule a number of tasks
	for i := 0; i < 10; i++ {
		n.Go(newTask(i))
	}

	// wait for all tasks to complete
	err := n.Wait()
	if err != nil {
		fmt.Printf("There was an error: %v", err)
		return
	}

	fmt.Println("Done!")
}

func newTask(id int) nursery.Task {

	return func(ctx context.Context) error {
		fmt.Printf("Task %d starts\n", id)

		// do heaby work
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)

		fmt.Printf("Task %d finished\n", id)
		return nil
	}
}
