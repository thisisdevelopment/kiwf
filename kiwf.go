// Package kiwf Kill It With Fire, is an advanced, yet simple internal heartbeat library for golang. The library removes the need to implement your own timeout procedures using ctx/timeout/deadline/cancel routines. Timeout, delay and startup time can be configured. Also there's a possibility to pass in a map[string]interface{} that will be passed into a configurable ExitFunction. If no ExitFunction is configured a panic will be generated.
package kiwf

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Kiwf struct {
	ticked time.Time
	mtx    *sync.RWMutex
	wg     *sync.WaitGroup
	cfg    *Config
	title  string
	closer context.CancelFunc
	ctx    context.Context
}

// New instantiates a new kiwf, title required, config can be nil. Check defaults
func New(title string, config *Config) (*Kiwf, error) {
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	if config == nil {
		config = new(Config)
	}
	config.setDefaults()
	ctx, cancel := context.WithCancel(context.Background())

	return &Kiwf{
		ticked: time.Now(),
		mtx:    &sync.RWMutex{},
		wg:     &sync.WaitGroup{},
		cfg:    config,
		title:  title,
		closer: cancel,
		ctx:    ctx,
	}, nil

}

// Start the waiting, this will be run in a go-routine
func (k *Kiwf) Start() {
	// wait the startup delay
	time.Sleep(k.cfg.DelayStartupTime)
	k.wg.Add(1)
	go func() {
		// at least sleep 1 janitor cycle
		time.Sleep(k.cfg.janitorInterval)
		k.Tick()

	LOOP:
		for {
			select {
			case <-k.ctx.Done():
				break LOOP
			default:
				k.mtx.RLock()
				if time.Since(k.ticked) < k.cfg.Timeout+k.cfg.janitorInterval {
					k.mtx.RUnlock()
					time.Sleep(k.cfg.janitorInterval)
					continue
				}
				k.mtx.RUnlock()
				if k.cfg.ExitFunction == nil {
					k.defaultExit(k.title, k.cfg.Passtru, time.Since(k.LastAction()))
				} else {
					fn := *k.cfg.ExitFunction
					fn(k.title, k.cfg.Passtru)
					// not sure if we paniced here, so just stop the ticker
					k.closer()
					break LOOP
				}
			}
		}
		k.wg.Done()
	}()
}

func (k *Kiwf) Close() {
	k.closer()
	k.wg.Wait()
}

func (k *Kiwf) Tick() bool {
	if !k.mtx.TryLock() {
		return false
	}
	defer k.mtx.Unlock()
	k.ticked = time.Now()
	return true
}

// Lastaction returns the last touched time
func (k *Kiwf) LastAction() time.Time {
	k.mtx.RLock()
	defer k.mtx.RUnlock()
	return k.ticked
}

func (k *Kiwf) defaultExit(title string, passtru map[string]interface{}, lastSince time.Duration) {
	panic(
		fmt.Sprintf("Killed it with fire '%s' time expired last action %v ago. set timeout %v. passtru vars %v, time obj %+v",
			title,
			lastSince,
			k.cfg.Timeout,
			passtru,
			k.ticked,
		))
}
