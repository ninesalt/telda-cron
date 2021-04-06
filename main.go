package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID        string
	Timeout   time.Duration
	Frequency time.Duration
	Function  func(ctx context.Context, done chan int)
}

// This function just prints OP1 and then after 6 seconds prints
// OP2, the delay is used to test the timeout functionality
func withDelay(ctx context.Context, done chan int) {
	// the following goroutine is just to simulate an action
	// that takes some time
	channel := make(chan time.Time)

	go func() {
		time.Sleep(6 * time.Second)
		channel <- time.Now()
	}()

	log.Println("OP1")
	select {
	case <-channel:
		// indicates the job has completed
		log.Println("OP2")
		done <- 0
	case <-ctx.Done():
		// indicates a timeout
		done <- 1
	}
}

// a regular function that runs without a delay
func withoutDelay(ctx context.Context, done chan int) {
	log.Println("Running without a delay")
}

func New(ID, frequency, runtime string,
	implementation func(ctx context.Context, done chan int)) Job {
	r, err := time.ParseDuration(runtime)

	if err != nil {
		panic(err)
	}

	f, err := time.ParseDuration(frequency)

	if err != nil {
		panic(err)
	}

	j := Job{ID: ID, Timeout: r, Frequency: f, Function: implementation}
	log.Printf("Created job %#v with frequency %v and timeout %v", ID, f, r)
	return j
}

func (j Job) Run() {
	for range time.Tick(j.Frequency) {
		instanceID := rand.Intn(10000)
		start := time.Now()
		log.Printf("Job %#v IID=%v executing...", j.ID, instanceID)

		// run in a goroutine so it doesnt block the ticker
		go func() {
			done := make(chan int)
			ctx, cancel := context.WithTimeout(context.Background(), j.Timeout)
			defer cancel()
			go j.Function(ctx, done)

			// zero exit code in channel indicates success, otherwise it's a timeout
			if <-done == 0 {
				elapsed := time.Since(start)
				log.Printf("Job %#v IID=%v completed in %v", j.ID, instanceID, elapsed)
			} else {
				log.Printf("Job %#v IID=%v exited due to timeout", j.ID, instanceID)
			}
		}()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// create a new job given its name, frequency, timeout
	// and the function it should run
	job1 := New("my-first-job", "1s", "5s", func(ctx context.Context, done chan int) {
		withDelay(ctx, done)
	})

	job2 := New("my-second-job", "3s", "10s", func(ctx context.Context, done chan int) {
		withoutDelay(ctx, done)
	})

	// run the jobs concurrently and wait for them using a wait group
	var wg sync.WaitGroup
	wg.Add(2)

	go job1.Run()
	go job2.Run()
	wg.Wait()
}
