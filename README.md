## Telda Cron


This project aims to simulate the basic functionality of cron by allowing users to create jobs with a frequency, timeout and a function which will be called at every tick of the given frequency. Supports regular frequency and timeout are parsed using `time.Duration` so the following are valid: `1h24m`, `"10s"`, `"200ms"`.

The implementation function which will be called should be passed a context object and a `done` channel. The context is used to check when the timeout occurs in which case a non-zero integer should be sent through the done channel to signal a timeout. A zero code in the done channel signifies a successful run of the job.


### Assumptions
- Using `rand.Intn(10000)` to generate instance IDs of running jobs (this is different than the ID the user passes when creating a job). Ideally uuids would be used but I stuck to this just out of simplicity and to avoid external deps.
- The timeout is required
