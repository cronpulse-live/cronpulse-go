package cronpulse

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

var exit = os.Exit // Create a variable to hold os.Exit so we can override it in tests

// Monitor represents the monitoring object.
type Monitor struct {
	jobKey string
}

// NewMonitor creates a new Monitor instance with the hardcoded base URL.
func NewMonitor(jobKey string) *Monitor {
	return &Monitor{jobKey: jobKey}
}

// Ping sends a ping request to the monitoring server.
func (m *Monitor) Ping(stateOrOptions interface{}) error {
	var state, message string

	switch v := stateOrOptions.(type) {
	case string:
		state = v
	case map[string]string:
		state = v["state"]
		message = v["message"]
	default:
		return fmt.Errorf("invalid argument: stateOrOptions must be a string or a map[string]string")
	}

	queryParams := url.Values{}
	queryParams.Set("client", "cronpulse go")

	var endpoint string
	switch state {
	case "beat":
		endpoint = fmt.Sprintf("/api/ping/%s", m.jobKey)
	case "start":
		endpoint = fmt.Sprintf("/api/ping/%s/start", m.jobKey)
	case "success":
		endpoint = fmt.Sprintf("/api/ping/%s/success", m.jobKey)
	case "fail":
		endpoint = fmt.Sprintf("/api/ping/%s/fail", m.jobKey)
		if message != "" {
			queryParams.Set("errorMessage", message)
		} else {
			queryParams.Set("errorMessage", "true")
		}
	default:
		return fmt.Errorf("invalid state: %s", state)
	}

	return m.sendRequest(endpoint, queryParams)
}

func (m *Monitor) sendRequest(endpoint string, queryParams url.Values) error {
	baseURL := "https://app.cronpulse.live"
	url := fmt.Sprintf("%s%s?%s", baseURL, endpoint, queryParams.Encode())
	fmt.Printf("🏓 Pinging : %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ Error: %s\n", err.Error())
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("📊 Status Code: %d\n", resp.StatusCode)
	return nil
}

// Wrap wraps a job function with monitoring.
func Wrap(jobKey string, jobFunction func() error) func() {
	monitor := NewMonitor(jobKey)
	startTime := time.Now()
	exitCode := 0
	errorOccurred := false
	errorMessage := ""

	originalExit := exit
	exitCalled := false

	exit = func(code int) {
		if exitCalled {
			return
		}
		exitCalled = true
		exitCode = code

		finalState := "success"
		if errorOccurred {
			finalState = "fail"
		}
		monitor.Ping(map[string]string{"state": finalState, "message": errorMessage})
		originalExit(exitCode)
	}

	return func() {
		defer func() {
			endTime := time.Now()
			fmt.Printf("Job execution time: %d ms\n", endTime.Sub(startTime).Milliseconds())

			if !exitCalled {
				finalState := "success"
				if errorOccurred {
					finalState = "fail"
				}
				monitor.Ping(map[string]string{"state": finalState, "message": errorMessage})
				exit(exitCode)
			}
		}()

		err := monitor.Ping(map[string]string{"state": "start"})
		if err != nil {
			errorOccurred = true
			errorMessage = err.Error()
			exitCode = 1
			return
		}

		err = jobFunction()
		if err != nil {
			errorOccurred = true
			errorMessage = err.Error()
			exitCode = 1
		}
	}
}
