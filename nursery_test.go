package nursery

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func ExampleNew() {
	n := New(context.Background())

	n.Go(func(ctx context.Context) error {
		time.Sleep(20 * time.Millisecond)
		fmt.Println("1 done")
		return errors.New("error")
	})
	n.Go(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("2 done")
		return nil
	})
	err := n.Wait()

	fmt.Printf("%v\n", err)

}

func TestWaitGroup(t *testing.T) {
	n := New(context.Background())

	complete1 := false
	complete2 := false
	n.Go(func(ctx context.Context) error {
		time.Sleep(20 * time.Millisecond)
		complete1 = true
		return nil
	})
	n.Go(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		complete2 = true
		return nil
	})
	err := n.Wait()

	if err != nil {
		t.Fail()
	}

	if !complete1 {
		t.Error("Task 1 did not complete")
	}
	if !complete2 {
		t.Error("Task 2 did not complete")
	}

}

func TestProducerWorkerConsumer(t *testing.T) {
	n := New(context.Background())

	prodChan := make(chan int)
	workerChan := make(chan string)
	n.Go(func(ctx context.Context) error {
		//Producer
		defer close(prodChan)
		for i := 0; i < 10; i++ {
			prodChan <- i
		}
		return nil
	})
	n.Go(func(ctx context.Context) error {
		//worker
		defer close(workerChan)
		for x := range prodChan {
			workerChan <- "as string : " + strconv.Itoa(x)
			time.Sleep(10 * time.Millisecond)
		}
		return nil
	})

	n.Go(func(ctx context.Context) error {
		// consumer
		for s := range workerChan {
			t.Log("consumer ", s)
		}
		return nil
	})
	err := n.Wait()
	if err != nil {
		t.Errorf("Wait got error: %v", err)
	}
}

func TestProducerWorkersConsumer(t *testing.T) {
	n := New(context.Background())

	prodChan := make(chan int)
	workerChan := make(chan string)
	n.Go(func(ctx context.Context) error {
		//Producer
		defer close(prodChan)
		for i := 0; i < 10; i++ {
			prodChan <- i
		}
		return nil
	})
	n.Go(func(ctx context.Context) error {
		//worker manager
		defer close(workerChan)
		wn := New(ctx)
		for j := 0; j < 2; j++ {
			id := strconv.Itoa(j)
			wn.Go(func(ctx context.Context) error {
				//worker
				for x := range prodChan {
					workerChan <- id + ": " + strconv.Itoa(x)
					time.Sleep(10 * time.Millisecond)
				}
				return nil
			})
		}
		wn.Wait()

		return nil
	})

	n.Go(func(ctx context.Context) error {
		// consumer
		for s := range workerChan {
			t.Log("consumer ", s)
		}
		return nil
	})
	err := n.Wait()
	if err != nil {
		t.Errorf("Wait got error: %v", err)
	}
}

func TestCancel(t *testing.T) {
	n := New(context.Background())

	for i := 0; i < 10; i++ {
		n.Go(func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					return nil
				default:
				}
				time.Sleep(10 * time.Millisecond)
			}
		})
	}
	n.Cancel()
	err := n.Wait()

	if err != ErrContextCanceled {
		t.Errorf("expected ErrTimeout, got :%v", err)
	}

}

func TestContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	n := New(ctx)

	for i := 0; i < 10; i++ {
		n.Go(func(ctx context.Context) error {
			for {
				time.Sleep(10 * time.Millisecond)
			}
		})
	}
	cancel()
	err := n.Wait()

	if err != ErrContextCanceled {
		t.Errorf("expected ErrTimeout, got :%v", err)
	}

}
