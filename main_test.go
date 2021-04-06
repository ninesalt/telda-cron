package main

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"
)

func testWithoutDelay(ctx context.Context, done chan int) {
	fmt.Println("hello")
}

func testWithDelay(ctx context.Context, done chan int, ioChannel chan string) {
	defer close(ioChannel)
	defer close(done)

	channel := make(chan time.Time)

	go func() {
		defer close(channel)
		time.Sleep(5 * time.Second)
		channel <- time.Now()
	}()

	ioChannel <- "OP1"
	select {
	case <-channel:
		ioChannel <- "OP2"
		done <- 0
	case <-ctx.Done():
		done <- 1
	}
}

func TestCreateJob(t *testing.T) {
	defer func() { recover() }()

	var f = func(ctx context.Context, done chan int) {
		testWithoutDelay(ctx, done)
	}

	New("test-job", "invalid", "1m", f)
	t.Errorf("Invalid frequency should trigger panic")

	New("test-job", "10s", "invalid", f)
	t.Errorf("Invalid timeout should trigger panic")
}

func TestRunJobWithoutTimeout(t *testing.T) {
	var ioChannel = make(chan string)
	var testDone = make(chan int)

	var f = func(ctx context.Context, done chan int) {
		testWithDelay(ctx, testDone, ioChannel)
	}

	last := time.Now()
	go New("WITHOUT_TIMEOUT-", "2s", "10s", f).Run()

	// listen on 3 messages coming in through the ioChannel
	// if its OP1, make sure the tick is 2 seconds
	for i := 0; i < 3; i, last = i+1, time.Now() {

		// this job should never trigger a timeout because
		// it takes 2 seconds to run and the timeout here is 10s
		value := <-ioChannel
		diff := math.Round(time.Since(last).Seconds())

		if value == "OP1" && diff != 2 {
			t.Error("Job should run exactly at every specified tick (2s)")
			return
		}
	}
}

func TestRunJobWithTimeout(t *testing.T) {
	var ioChannel = make(chan string)
	var testDone = make(chan int)

	var f = func(ctx context.Context, done chan int) {
		testWithDelay(ctx, testDone, ioChannel)
	}

	// this job will sleep for 5 seconds and it has a timeout of 4 seconds
	// so it should timeout and exit
	go New("WITH_TIMEOUT", "2s", "4s", f).Run()

	begunJobs := 0
	hasTimedOut := false

L:
	for begunJobs < 4 {
		// this job should never trigger a timeout because
		// it takes 2 seconds to run and the timeout here is 4s
		// the job takes 5 seconds to run so it should timeout exactly
		// as the second job instance starts
		select {
		case value := <-ioChannel:
			if value == "OP1" {
				begunJobs++
			} else { // if we get OP2 then the job completed which is wrong
				break L
			}
		case doneValue := <-testDone:
			if doneValue == 1 {
				hasTimedOut = true
				break L // exit out of the loop
			}
		}
	}

	if hasTimedOut == false {
		t.Errorf("Job should timeout after specified timeout")
	}

}
