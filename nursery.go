package nursery

import (
	"context"
	"errors"
	"sync"
)

type Task func(context.Context) error

var ErrTimeout = errors.New("Timeout")
var ErrContextCanceled = errors.New("Context Canceled")

type Nursery struct {
	errChan chan error
	ctx     context.Context
	cancel  func()
	wg      sync.WaitGroup
}

func New(ctx context.Context) *Nursery {
	ctx, cancel := context.WithCancel(ctx)

	return &Nursery{
		ctx:     ctx,
		cancel:  cancel,
		errChan: make(chan error),
	}
}

func (n *Nursery) Go(task Task) *Nursery {

	n.wg.Add(1)
	go func() {
		err := task(n.ctx)
		if err != nil {
			n.errChan <- err
		}
		n.wg.Done()
	}()
	return n
}

// Wait blocks until either one task has exited with an error
// or until all tasks have been finished without.
// If an error occures this error is retruned by Wait()
func (n *Nursery) Wait() error {
	doneCh := make(chan struct{})
	go func() {
		n.wg.Wait()
		close(doneCh)
	}()
	var err error
	select {
	case <-doneCh:
	case err = <-n.errChan:
		n.cancel()
	case <-n.ctx.Done():
		err = ErrContextCanceled
	}
	return err
}

// Cancel cancels the context provided to the tasks.
// You may want to Wait() to be sure that all tasks have exited.
func (n *Nursery) Cancel() {
	n.cancel()
}
