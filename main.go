package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Job struct {
	ID        string
	Timeout   time.Duration
	Frequency time.Duration
	Function  func(ctx context.Context)
}

func testFunc(ctx context.Context) {
	channel := make(chan time.Time)

	go func() {
		time.Sleep(7 * time.Second)
		channel <- time.Now()
	}()

	log.Println("OP1")
	select {
	case <-channel:
		fmt.Println("OP2")
	case <-ctx.Done():
		log.Println("Job exiting due to timeout")
		return
	}

}

func New(ID, frequency, runtime string, implementation func(ctx context.Context)) Job {
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
		// start := time.Now()
		log.Printf("Job %#v IID=%v executing...", j.ID, instanceID)

		// done := make(chan int, 1)
		// quit := make(chan int, 1)

		ctx, cancel := context.WithTimeout(context.Background(), j.Timeout)
		defer cancel()

		// go func() {
		go j.Function(ctx)
		// done <- 0
		// select {
		// case <-done:
		// 	elapsed := time.Since(start)
		// 	log.Printf("Job %#v IID=%v completed in %v \n", j.ID, instanceID, elapsed)
		// }
		// }()

	}

}

func main() {
	rand.Seed(time.Now().UnixNano())

	// create a new job given its name, frequency, max runtime
	// and the function it should run
	testJob := New("my-first-job", "3s", "5s", func(ctx context.Context) {
		testFunc(ctx)
	})

	testJob.Run()
}
