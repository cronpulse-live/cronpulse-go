# cronpulse-go

`cronpulse-go` is a simple Go package for monitoring job execution using Cronpulse.live. This package allows you to send heartbeats, start, success, and failure signals to monitor the status of your cron jobs or any scheduled tasks.

## Installation

To install the package, run:

```sh
go get github.com/cronpulse-live/cronpulse-go
```

## Usage

Here's an example of how to use `cronpulse-go`:

```go
package main

import (
    "github.com/cronpulse-live/cronpulse-go"
    "fmt"
)

func main() {
    monitor := cronpulse.NewMonitor("your-job-key")

    err := monitor.Ping("start")
    if err != nil {
        fmt.Printf("Error pinging start: %v\n", err)
    }

    // Your job logic here...

    err = monitor.Ping("success")
    if err != nil {
        fmt.Printf("Error pinging success: %v\n", err)
    }
}
```

### Using the Wrap Function

The `Wrap` function can be used to automatically monitor the job's start, success, and failure status.

```go
package main

import (
    "github.com/cronpulse-live/cronpulse-go"
    "fmt"
)

func main() {
    job := cronpulse.Wrap("your-job-key", func() error {
        // Your job logic here...
        return nil // or return an error if the job fails
    })

    job()
}
```

## Testing

To run the tests:

```sh
go test -v
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

MIT License

Copyright (c) 2024 Your Name

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT
