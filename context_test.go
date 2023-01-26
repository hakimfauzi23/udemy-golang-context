package golang_context_udemy

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

/* Create New Empty Context */
func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

/* Create Context With Value And Implement Child-Parent Context */
func TestContextWithValue(t *testing.T) {
	contextA := context.Background()
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)

	fmt.Println(contextF.Value("f")) // get value with key 'f'
	fmt.Println(contextF.Value("c"))
	fmt.Println(contextF.Value("b")) // nil because not related.

	fmt.Println(contextA.Value("b")) // nil because context is up direction (child to parent).

}

// func CreateCounter() chan int {
// 	destination := make(chan int)

// 	go func() {
// 		defer close(destination)
// 		counter := 1
// 		for {
// 			destination <- counter
// 			counter++
// 		}
// 	}()

//		return destination
//	}

func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
			}
		}
	}()

	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()

	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}

	// Sending cancel signal to context to avoid goroutine leak
	cancel()

	time.Sleep(2 * time.Second)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

// Goroutine Leak => Goroutine always running even the program is not need it.
// Much Goroutine Leak will cause slows and even failure stops
