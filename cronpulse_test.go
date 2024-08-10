package cronpulse

import (
	"errors"
	"fmt"
	"testing"
)

// Mock exit function
var testExitCalled = false
var testExitCode int

func mockExit(code int) {
	testExitCalled = true
	testExitCode = code
}

// TestMonitorPingSuccess tests the Monitor's Ping function for a start -> success sequence.
func TestMonitorPingSuccess(t *testing.T) {
	monitor := NewMonitor("ft1joc4eo")

	t.Log("🚀 Testing start -> success scenario")

	// Send start ping
	if err := monitor.Ping("start"); err != nil {
		t.Fatalf("❌ Error during start ping: %v", err)
	}

	// Send success ping
	if err := monitor.Ping("success"); err != nil {
		t.Fatalf("❌ Error during success ping: %v", err)
	}

	t.Log("✅ Successfully completed start -> success scenario")
}

// TestMonitorPingFailure tests the Monitor's Ping function for a start -> fail sequence.
func TestMonitorPingFailure(t *testing.T) {
	monitor := NewMonitor("ft1joc4eo")

	t.Log("🚀 Testing start -> failure scenario")

	// Send start ping
	if err := monitor.Ping("start"); err != nil {
		t.Fatalf("❌ Error during start ping: %v", err)
	}

	// Create a dynamic error message
	dynamicError := errors.New("This is a dynamic error message")

	// Send failure ping with dynamic error message
	if err := monitor.Ping(map[string]string{"state": "fail", "message": dynamicError.Error()}); err != nil {
		t.Fatalf("❌ Error during failure ping: %v", err)
	}

	t.Log("✅ Successfully completed start -> failure scenario")
}

// TestMonitorPingBeat tests the Monitor's Ping function for a beat (heartbeat) state.
func TestMonitorPingBeat(t *testing.T) {
	t.Log("🚀 Testing beat (heartbeat) scenario")
	monitor := NewMonitor("ft1joc4eo")

	if err := monitor.Ping("beat"); err != nil {
		t.Fatalf("❌ Error during beat ping: %v", err)
	}

	t.Log("✅ Successfully completed beat (heartbeat) scenario")
}

// TestWrap tests the Wrap function with both success and failure scenarios.
func TestWrap(t *testing.T) {
	monitor := NewMonitor("ft1joc4eo")

	// Mock exit function
	originalExit := exit
	exit = mockExit
	defer func() { exit = originalExit }()

	// Success case
	t.Log("🚀 Testing wrap function with success scenario")
	jobFuncSuccess := func() error {
		return nil
	}

	testExitCalled = false
	wrappedSuccess := Wrap(monitor.jobKey, jobFuncSuccess)
	wrappedSuccess()

	if testExitCalled {
		t.Errorf("❌ Expected exit not to be called, but it was")
	} else {
		t.Log("✅ Successfully completed wrap function with success scenario")
	}

	// Failure case
	t.Log("🚀 Testing wrap function with failure scenario")
	dynamicError := fmt.Errorf("Dynamic error: %s", "something went wrong")
	jobFuncFailure := func() error {
		return dynamicError
	}

	testExitCalled = false
	testExitCode = 0
	wrappedFailure := Wrap(monitor.jobKey, jobFuncFailure)
	wrappedFailure()

	if !testExitCalled {
		t.Errorf("❌ Expected exit to be called, but it wasn't")
	} else if testExitCode != 1 {
		t.Errorf("❌ Expected exit code 1, but got %d", testExitCode)
	} else {
		t.Log("✅ Successfully completed wrap function with failure scenario")
	}
}

// TestMonitorPingInvalidState tests the Monitor's Ping function with an invalid state.
func TestMonitorPingInvalidState(t *testing.T) {
	t.Log("🚀 Testing invalid state scenario")
	monitor := NewMonitor("ft1joc4eo")

	if err := monitor.Ping("invalidState"); err == nil {
		t.Errorf("❌ Expected an error for invalid state, got none")
	} else {
		t.Logf("✅ Received expected error for invalid state: %v", err)
	}
}
