## Telda Cron


This project aims to simulate the basic functionality of cron by allowing users to create jobs with a frequency, timeout and a function which will be called at every tick of the given frequency. Supports regular frequency and timeout are parsed using `time.Duration` so the following are valid: `"1h24m"`, `"10s"`, `"200ms"`.

The implementation function which will be called should be passed a context object and a `done` channel. The context is used to check when the timeout occurs in which case a non-zero integer should be sent through the done channel to signal a timeout. A zero code in the done channel signifies a successful run of the job.

## Example
```
func myFunc(ctx context.Context, done chan int, name string) {
	log.Printf("Hello, %v!\n", name)
	done <- 0
}

func main() {
	testJob := New("my-job", "5s", "30s", func(ctx context.Context, done chan int) {
		myFunc(ctx, done, "youssef")
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go testJob.Run()
    
    // run more jobs concurrently here

    defer testJob.Stop()
	wg.Wait()
}
```

This should output:
```
2021/04/07 12:24:58 Created job "my-job" with frequency 5s and timeout 30s
2021/04/07 12:25:03 Job "my-job" IID=8081 executing...
2021/04/07 12:25:03 Hello, youssef!
2021/04/07 12:25:03 Job "my-job" IID=8081 completed in 173.729µs
2021/04/07 12:25:08 Job "my-job" IID=7887 executing...
2021/04/07 12:25:08 Hello, youssef!
2021/04/07 12:25:08 Job "my-job" IID=7887 completed in 109.506µs
2021/04/07 12:25:13 Job "my-job" IID=1847 executing...
2021/04/07 12:25:13 Hello, youssef!
2021/04/07 12:25:13 Job "my-job" IID=1847 completed in 142.848µs
```
### Assumptions
- Using `rand.Intn(10000)` to generate instance IDs of running jobs (this is different than the ID the user passes when creating a job). I added IIDs so I can more easily debug log messages to figure out which instance logged which message at what time. Ideally uuids would be used but I stuck to this just out of simplicity and to avoid external deps.
- The timeout is required