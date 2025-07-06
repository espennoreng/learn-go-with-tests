package main

import (
	"os"
	"time"

	"github.com/espennoreng/learn-go-with-tests/mocking"
)

func main() {
	sleeper := &mocking.ConfigurableSleeper{
		Duration:  1 * time.Second,
		SleepFunc: time.Sleep,
	}
	mocking.Countdown(os.Stdout, sleeper)
}
