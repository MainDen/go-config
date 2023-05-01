# Go-Config

Flexible developer oriented configuration library.

[![ci status][ci-status-badge]][ci-status-link]
[![code coverage][codecov-badge]][codecov-link]
[![go report][goreportcard-badge]][goreportcard-link]
[![go dev][go-dev-badge]][go-dev-link]
[![repo size][repo-size-badge]][repo-size-link]
[![donation][paypal-badge]][paypal-link]
[![license][license-badge]][license-link]

[ci-status-badge]: https://img.shields.io/github/actions/workflow/status/MainDen/go-config/tests.yml?branch=main&event=push&logo=github&label=tests
[ci-status-link]: https://github.com/MainDen/go-config/actions/workflows/tests.yml?query=branch%3Amain+event%3Apush
[codecov-badge]: https://img.shields.io/codecov/c/github/mainden/go-config/main?token=NAK406E7A6&logo=codecov
[codecov-link]: https://codecov.io/gh/MainDen/go-config/tree/main
[goreportcard-badge]: https://goreportcard.com/badge/github.com/MainDen/go-config
[goreportcard-link]: https://goreportcard.com/report/github.com/MainDen/go-config
[go-dev-badge]: https://img.shields.io/badge/go-reference-blue?logo=go&logoColor=white
[go-dev-link]: https://pkg.go.dev/github.com/MainDen/go-config
[repo-size-badge]: https://img.shields.io/github/repo-size/MainDen/go-config?logo=github
[repo-size-link]: https://github.com/MainDen/go-config/tree/main
[paypal-badge]: https://img.shields.io/badge/donate-paypal-blue?logo=paypal
[paypal-link]: https://www.paypal.me/mainden
[license-badge]: https://img.shields.io/badge/license-MIT-yellow
[license-link]: https://github.com/MainDen/go-config/blob/main/LICENSE.md

# Overview

- Pointers configuration
- Type convertion
- Detailed errors
- Logging

# Install

```cmd
go get github.com/mainden/go-config
```

# Quick Start

```go
package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mainden/go-config/configurator"
)

type AppConfig struct {
	UserName string
	Token    string
	Timeout  time.Duration
}

func main() {
	// Create new configurator or use configurator.Default.
	configurator := configurator.NewConfigurator().WithLogger(log.Printf)

	config := AppConfig{}
	// UserName should not be empty. Use 'root' by default.
	err := configurator.WithDisallowed("").WithDefault("root").Configure(&config.UserName)
	if err != nil {
		log.Fatal(err)
	}
	// Token should not be empty. Don't log token value (secret).
	err = configurator.WithDisallowed("").WithLogSecret().Configure(&config.Token)
	if err != nil {
		log.Fatal(err)
	}
	// Timeout should not be lower than second and should not be greater than minute.
	err = configurator.WithMin(time.Second).WithMax(time.Minute).Configure(&config.Timeout)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Timeout: config.Timeout,
	}
	url, err := url.Parse("https://example.com/")
	if err != nil {
		log.Fatal(err)
	}
	request := &http.Request{
		URL:    url,
		Method: http.MethodPost,
		Header: http.Header{
			"UserName": {config.UserName},
			"Token":    {config.Token},
		},
		Body: io.NopCloser(strings.NewReader("Some Data")),
	}
	_, err = client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
}
```
