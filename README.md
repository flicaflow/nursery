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
    #--------+   +--------+   +--------+
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


