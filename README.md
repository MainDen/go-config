# Go-Config

Flexible configuration package. Validate configuration with simple pattern.

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
[go-dev-badge]: https://img.shields.io/badge/go%20doc-reference-blue?logo=go&logoColor=white
[go-dev-link]: https://pkg.go.dev/github.com/mainden/go-config
[repo-size-badge]: https://img.shields.io/github/repo-size/MainDen/go-config?logo=github
[repo-size-link]: https://github.com/MainDen/go-config/tree/main
[paypal-badge]: https://img.shields.io/badge/donate-paypal-blue?logo=paypal
[paypal-link]: https://www.paypal.me/mainden
[license-badge]: https://img.shields.io/badge/license-MIT-yellow
[license-link]: https://github.com/MainDen/go-config/blob/main/LICENSE.md

# Overview

A lot of features are supported:
- Pointers configuration (compare values instead of pointers)
- Type convertion (don't care about derived types)
- Detailed errors (you see a field that is invalid and why)
- Logging (provide your custom or default logger)

# Install

```cmd
go get github.com/mainden/go-config/...
```

# Quick Start

Basically we have application config and a lot of validations. We can sipmlify a lot of lines of code with simple chained call. You don't need define many validation methods for every type, just use configuring package!

```go
// AppConfig is our application config.
type AppConfig struct {
	Address net.IP
	Timeout time.Duration
}

// Configure is method to configure our application config.
func Configure(config *AppConfig) error {
	// We can define a default Address that should not be empty. Let's use the '127.0.0.1' by default for example.
	if err := configuring.Default.WithName("Address").WithDisallowed(net.IP(nil)).WithDefault(net.IPv4(127, 0, 0, 1)).Configure(&config.Address); err != nil {
		return err
	}
	// Let's configure Timeout that should not be lower than second and should not be greater than minute.
	if err := configuring.Default.WithName("Timeout").WithMin(time.Second).WithMax(time.Minute).Configure(&config.Timeout); err != nil {
		return err
	}
	return nil
}
```

Now we can configure 'AppConfig' with the 'Configure' method.

```go
	var config AppConfig

	// Configure(&config) => error: "configuration of 'Timeout' error: target value error: argument should be greater than or equal to '1s'"

	config.Timeout = 2 * time.Second
	// Value should be set to prevent previous error

	Configure(&config) // no error
	// Here we have 'config' that equal to AppConfig{Address: net.IPv4(127, 0, 0, 1), Timeout: 2 * time.Second}
	// Log:
	// configuration of 'Address': disallowed: ['<nil>'] default: '127.0.0.1' input: '<nil>' output: '127.0.0.1'
	// configuration of 'Timeout': min: '1s' max: '1m0s' input: '2s' output: '2s'
```
