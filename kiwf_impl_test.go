package kiwf

import (
	"testing"
	"time"
)

func Test_kiwf_Impl_1(t *testing.T) {

	var (
		timeout    = 1 * time.Second
		exitCalled = false
		exitFunc   = func(_ string, _ map[string]interface{}) {
			exitCalled = true
		}
	)

	kiwf, err := New("test case 1", &Config{
		DelayStartupTime: 50 * time.Millisecond,
		Timeout:          timeout,
		ExitFunction:     &exitFunc,
	})
	if err != nil {
		t.FailNow()
	}

	kiwf.Start()
	defer kiwf.Close()
	time.Sleep(200 * time.Millisecond)

	if time.Since(kiwf.LastAction()) > timeout {
		t.Errorf("should be longer than timeout,last touched %v, got timeout %v", time.Since(kiwf.LastAction()), timeout)
	}

	if exitCalled {
		t.Errorf("exit called but no timeout occurred, last action %v", time.Since(kiwf.LastAction()))
	}

}

func Test_kiwf_Impl_2(t *testing.T) {
	var exitCalled = false
	var testTitle = ""

	var (
		endTime = time.Now()
		timeout = 1 * time.Second
	)

	exitFunc := func(title string, _ map[string]interface{}) {
		testTitle = title
		endTime = time.Now()
		exitCalled = true
	}

	tt := "t2"
	kiwf, err := New(tt, &Config{
		DelayStartupTime: 1 * time.Second,
		Timeout:          timeout,
		ExitFunction:     &exitFunc,
	})
	if err != nil {
		t.FailNow()
	}

	kiwf.Start()
	defer kiwf.Close()

	time.Sleep(2200 * time.Millisecond)

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

	kiwf, err := New("tc3", &Config{
		DelayStartupTime: time.Nanosecond,
		Timeout:          5 * time.Millisecond,
	})
	if err != nil {
		t.FailNow()
	}

	startTime := time.Now()
	kiwf.Start()
	defer kiwf.Close()

	for time.Since(startTime) < 50*time.Millisecond {
		kiwf.Tick()
		time.Sleep(time.Millisecond)
	}
}

func Test_kiwf_Impl_4(t *testing.T) {

	var (
		testTimeout = 500 * time.Millisecond
	)
	_, err := New("", &Config{})
	if err == nil {
		t.Error("should throw error empty title")
	}
	kiwf, err := New("test case nil config", nil)
	if err != nil {
		t.FailNow()
	}

	kiwf.Start()
	defer kiwf.Close()
	startTime := time.Now()
	for time.Since(startTime) < testTimeout {
		kiwf.Tick()
		time.Sleep(1 * time.Millisecond)
	}

}

func Test_kiwf_Impl_5(t *testing.T) {
	// passtru test
	var (
		exitCalled = false
		testTitle  = ""
		testVar    = ""
		testPassed = "test123"
		endTime    = time.Now()
		timeout    = 1 * time.Second
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
	defer kiwf.Close()
	time.Sleep(2200 * time.Millisecond)

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

	// simulation settings
	var (
		testTimeout = 500 * time.Millisecond
	)

	kiwf, err := New("my work", &Config{
		DelayStartupTime: 10 * time.Millisecond,
		Timeout:          150 * time.Millisecond,
	})
	if err != nil {
		t.Errorf("got error '%s'", err.Error())
	}
	startTime := time.Now()
	kiwf.Start()
	defer kiwf.Close()

	// simulate work and tick in iteration
	for time.Since(startTime) < testTimeout {
		kiwf.Tick()
		time.Sleep(1 * time.Millisecond)
	}
	time.Sleep(time.Second)
	// yep this panics
}
**/
