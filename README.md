# README IS WIP 

[![go report card](https://goreportcard.com/badge/github.com/thisisdevelopment/kiwf "go report card")](https://goreportcard.com/report/github.com/thisisdevelopment/kiwf)


# Introduction
**Kill It With Fire** 🔥 🔫 is an advanced, yet simple internal heartbeat library for golang. The library removes the need to implement your own timeout procedures using ctx/timeout/deadline/cancel routines. Timeout, delay and startup time can be configured. Also there's a possibility to pass in a map[string]interface{} that will be passed into a configurable ExitFunction. If no ExitFunction is configured a panic will be generated.

# Basic usage 
```
    // initialize 
    kiwf, err := New("my work", &Config{
        // delay for sometime to start monitoring
		DelayStartupTime: 10 * time.Millisecond,

        // duration we should kill or call function if we didnt receive a Tick in time
		Timeout:          150 * time.Millisecond,
	})
	if err != nil {
		// do sth with err
	}


    // start watching 
	kiwf.Start()

    // always call Close
	defer kiwf.Close()

    testTimeout := 500 * time.Millisecond
	startTime := time.Now()

    // simulate some work for approx 500ms and timeout
	for {
		kiwf.Tick()
		if time.Since(startTime) > testTimeout {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}

    // simulate idleness, do nothing, take a nap
    time.Sleep(time.Second)

    // this should panic vvvv

    panic: Killed it with fire 'my work' time expired last action 150.833542ms ago. set timeout 150ms. passtru vars map[], time obj 2022-12-14 11:11:08.212809 +0100 CET m=+0.501677876


```


# Contributing 
You can help to deliver a better fire 🔥 🔫  killer, check out how you can do things [CONTRIBUTING.md](CONTRIBUTING.md)

# License 
© [This is Development BV](https://www.thisisdevelopment.nl), 2022~time.Now()
Released under the [MIT License](https://github.com/thisisdevelopment/fanunmarshal/blob/master/LICENSE)