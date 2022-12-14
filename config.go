package kiwf

import (
	"time"
)

const (
	defaultDelayStartupTime = 1 * time.Second
	defaultTimeout          = 1 * time.Second
	janitorFactor           = 10
)

type Config struct {
	// DelayStartupTime initial time the process waits to start the janitor, default is 1 second
	DelayStartupTime time.Duration

	// ExitFunction function to execute if timeout occurs, if nil a panic is invoked
	ExitFunction *func(title string, passtru map[string]interface{})

	// Timeout the duration considdered a timeout thus invoke the ExitFunction, default is 1 second
	Timeout time.Duration

	// Passtru any vars that are passed on into your exit function
	Passtru map[string]interface{}

	// janitorInterval check every x time duration, this should be private and is calculated from timeout settings
	janitorInterval time.Duration
}

// setDefaults set default values (see consts)
func (c *Config) setDefaults() {

	if c.Timeout.Nanoseconds() == 0 {
		c.Timeout = defaultTimeout
	}

	if c.janitorInterval.Nanoseconds() == 0 {
		c.janitorInterval = c.Timeout / janitorFactor
	}

	if c.DelayStartupTime.Nanoseconds() == 0 {
		c.DelayStartupTime = defaultDelayStartupTime
	}

}
