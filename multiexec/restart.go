package multiexec

import (
	"time"
)

// RestartPolicy is a function type to determine how long to wait between
// program exit and restart of the program. A return value of zero causes an
// immediate restart. A negative return value means the program is not
// restarted.
type RestartPolicy func(*Context) time.Duration

var (
	// RunOnce is a restart policy which causes the program to run only once.
	RunOnce = RestartAfterDelay(-1)

	// RestartImmediately is a restart policy which causes the program to
	// restart immediately.
	RestartImmediately = RestartAfterDelay(0)
)

// RestartAfterDelay returns a restart policy which restarts a program after
// the specified delay. If delay is zero, the program is not restarted.
func RestartAfterDelay(delay time.Duration) RestartPolicy {
	return func(*Context) time.Duration {
		return delay
	}
}

// RestartNTimes returns a restart policy which restarts a program the specified
// number of times, using the specified delay.
func RestartNTimes(n int, delay time.Duration) RestartPolicy {
	return func(*Context) time.Duration {
		if n > 0 {
			n--
			return delay
		}
		return -1
	}
}
