package mocking

import (
	"fmt"
	"io"
	"time"
)

type Sleeper interface {
	Sleep()
}

type DefaultSleeper struct{}

func (d *DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

const (
	finalWord      = "Go!"
	countdownStart = 3
)

func Countdown(writer io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(writer, i)
		sleeper.Sleep()
	}
	fmt.Fprint(writer, finalWord)
}

type ConfigurableSleeper struct {
	Duration  time.Duration
	SleepFunc func(time.Duration)
}

func (c ConfigurableSleeper) Sleep() {
	c.SleepFunc(c.Duration)
}
