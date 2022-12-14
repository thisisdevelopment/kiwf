package kiwf

import (
	"testing"
	"time"
)

var exitCalled = false
var testTitle = ""

func Test_kiwf_Impl_1(t *testing.T) {
	timeout := 1 * time.Second
	exitCalled = false

	exitFunc := func(_ string, _ map[string]interface{}) {
		exitCalled = true
	}

	kiwf, err := New("test case 1", &Config{
		DelayStartupTime: 50 * time.Millisecond,
		Timeout:          timeout,
		ExitFunction:     &exitFunc,
	})
	if err != nil {
		t.Errorf("got error '%s'", err.Error())
	}

	kiwf.Start()
	time.Sleep(200 * time.Millisecond)
	kiwf.Close()

	if time.Since(kiwf.LastAction()) > timeout {
		t.Errorf("should be longer than timeout,last touched %v, got timeout %v", time.Since(kiwf.LastAction()), timeout)
	}

	if exitCalled {
		t.Errorf("exit called but no timeout occurred, last action %v", time.Since(kiwf.LastAction()))
	}

}

func Test_kiwf_Impl_2(t *testing.T) {
	exitCalled = false
	testTitle = ""

	var (
		endTime = time.Now()
		timeout = 1 * time.Second
	)

	exitFunc := func(title string, _ map[string]interface{}) {
		testTitle = title
		endTime = time.Now()
		exitCalled = true
	}

	tt := "test case 2"
	kiwf, err := New(tt, &Config{
		DelayStartupTime: 1 * time.Second,
		Timeout:          timeout,
		ExitFunction:     &exitFunc,
	})
	if err != nil {
		t.Errorf("got error '%s'", err.Error())
	}
	kiwf.Start()
	time.Sleep(2200 * time.Millisecond)
	kiwf.Close()

	// this should timeout
	if time.Since(kiwf.LastAction()) < timeout {
		t.Errorf("should be longer than timeout,last touched %v, end function called %v, got timeout %v", time.Since(kiwf.LastAction()), time.Since(endTime), timeout)
	}

	if !exitCalled {
		t.Error("exit called should have been called")
	}

	if tt != testTitle {
		t.Errorf("test title should match title got '%s', expected '%s'", testTitle, tt)
	}

}

// Test_kiwf_Impl_3, we mimic a process that has a 5ms timeout running for 50ms and doing work
func Test_kiwf_Impl_3(t *testing.T) {

	// internal settings
	var (
		timeout          = 5 * time.Millisecond
		delayStartupTime = 1 * time.Nanosecond
	)

	// simulation settings
	var (
		testTimeout     = 50 * time.Millisecond
		workerSleepTime = 1 * time.Millisecond
	)

	tt := "test case 3"
	kiwf, err := New(tt, &Config{
		DelayStartupTime: delayStartupTime,
		Timeout:          timeout,
	})
	if err != nil {
		t.Errorf("got error '%s'", err.Error())
	}

	startTime := time.Now()
	kiwf.Start()
	defer kiwf.Close()

	for {
		kiwf.Tick()
		if time.Since(startTime) > testTimeout {
			break
		}
		time.Sleep(workerSleepTime)
	}
}

func Test_kiwf_Impl_4(t *testing.T) {
	exitCalled = false
	testTitle = ""
	var (
		testTimeout = 500 * time.Millisecond
	)
	_, err := New("", &Config{})
	if err == nil {
		t.Error("should throw error empty title")
	}

	kiwf, _ := New("test case nil config", nil)

	kiwf.Start()
	startTime := time.Now()
	for {
		kiwf.Tick()

		if time.Since(startTime) > testTimeout {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	kiwf.Close()
}

func Test_kiwf_Impl_5(t *testing.T) {
	// passtru test
	exitCalled = false
	testTitle = ""
	var testVar = ""
	var testPassed = "test123"
	var (
		endTime = time.Now()
		timeout = 1 * time.Second
	)

	exitFunc := func(title string, passtru map[string]interface{}) {
		testTitle = title
		endTime = time.Now()
		exitCalled = true
		testVar = passtru["test"].(string)
	}

	tt := "test case passtru"
	kiwf, _ := New(tt, &Config{
		DelayStartupTime: 1 * time.Second,
		Timeout:          timeout,
		ExitFunction:     &exitFunc,
		Passtru:          map[string]interface{}{"test": testPassed},
	})

	kiwf.Start()
	time.Sleep(2200 * time.Millisecond)
	kiwf.Close()

	// this should timeout
	if time.Since(kiwf.LastAction()) < timeout {
		t.Errorf("should be longer than timeout,last touched %v, end function called %v, got timeout %v", time.Since(kiwf.LastAction()), time.Since(endTime), timeout)
	}

	if !exitCalled {
		t.Error("exit called should have been called")
	}

	if tt != testTitle {
		t.Errorf("test title should match title got '%s', expected '%s'", testTitle, tt)
	}

	if testVar != testPassed {
		t.Errorf("passtru test failed, expected '%s', got '%s'", testPassed, testVar)
	}
}

/**
// example on md file
func Test_kiwf_Impl_6(t *testing.T) {

	// internal settings
	var (
		timeout          = 150 * time.Millisecond
		delayStartupTime = 10 * time.Millisecond
	)

	// simulation settings
	var (
		testTimeout     = 500 * time.Millisecond
		workerSleepTime = 1 * time.Millisecond
	)

	tt := "my work"
	kiwf, err := New(tt, &Config{
		DelayStartupTime: delayStartupTime,
		Timeout:          timeout,
	})
	if err != nil {
		t.Errorf("got error '%s'", err.Error())
	}
	startTime := time.Now()
	kiwf.Start()
	defer kiwf.Close()

	for {
		kiwf.Tick()
		if time.Since(startTime) > testTimeout {
			break
		}
		time.Sleep(workerSleepTime)
	}
	time.Sleep(time.Second)
	// yep this panics
}
**/
