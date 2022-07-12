[![Go Report Card](https://goreportcard.com/badge/github.com/flicaflow/nursery)](https://goreportcard.com/report/github.com/flicaflow/nursery) [![Go Reference](https://pkg.go.dev/badge/github.com/flicaflow/nursery.svg)](https://pkg.go.dev/github.com/flicaflow/nursery)


# Nursery: Simple Structered Concurency for Go

Nursery is a minimal library providing structured concurrency primitives for Go.
Go provides with the go key word basic facilities for coroutine bases concurrency. 
Synchronization is supported by classical primitives such as mutex or wait groups or more high level with channels.
Structuring the concurrent execution of work is up to the user.
This involves topics like waiting for goroutines to complete, and error handling.

Structured concurrency aims to give concurrent execution a structure which is easy to understand and reason about.
The most common structure of concurrent execution involves spawning of several goroutines, and waiting for then to finish.
```
                  +-------+
                  | Spawn |
                  +-------+
                      |
         +------------+-----------+
         |            |           |
    +--------+   +--------+   +--------+
    | Task 1 |   | Task 2 |   | Task 3 |
    +--------+   +--------+   +--------+
         |            |           |
         +------------+-----------+
                      |
                 +---------+
                 |   Wait  |
                 +---------+
```
The concurrent execution is considered to be successful when all Tasks finish without an error.

This library simplifies the implementation of this model by implementing a goroutine **nursery** which manages the execution synchronization and error reporting of the tasks.
```
	n := New(context.Background())
	n.Go(task1)
	n.Go(task2)
	n.Go(task3)
	err := n.Wait()
	if err!=nil {
		n.Cancel()
	}
```
Wait will block until either until either all task have finished without an error or until one task has finished with an error.

A task is a function with the signature `func(context.Context) error`. 
Receiving a context means that long running tasks can observe the `ctx.Done()` channel to terminate early. The context canceling  can done by calling `n.Cancel()` on the nursers and happens of course automatically when the provided context is canceled otherwise.

And that is basically everything. Three methods, `Go`, `Wait` and `Cancel` are enough to implement structured concurrency which greatly simplifies concurrent code. 

## Examples

### Simple task scheduling

Complete example in [examples/independent_tasks]
Complete example in [https://github.com/flicaflow/nursery/blob/main/examples/independent_tasks/independenttasks.go](examples/independent_tasks)

Nursery can be used to schedule and wait for any number of independent tasks.

```
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
```
`n.Wait()` will report any error returned by a task. It is possible to handle this errors in different ways.

* To stop the computation call `n.Cancel()`. You can wait again to observe finishing of remaining tasks and potential more errors.
* For error reporting, just wait in a loop until err is nil and log all errors.
* Conditional cancel, decide on the error type whether to cancel or continue.

### Compute pipelines

Complete example in [https://github.com/flicaflow/nursery/blob/main/examples/pipeline/pipeline.go](examples/pipeline)

More often tasks are not independent. Consider a simple processing pipeline with a producer, worker and collector communicating with channels.

```
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
```

In this case we cancel the nursery when an error occurs because we figure that the results collected by the collector are not complete and therefore worthless but that depends of course on the problem at hand.

The producer is implemented as follows:
```
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
```
It is observing the `ctx.Done()` channel. If the context is canceled this will exit the goroutine and close the work channel which will initiated the shutdown of the complete pipeline.
We didn't bother the handle context cancel in the other tasks since the channel close events will lead to finishing the pipeline anyway.
