# jsobfuscator-go

[![Go Report Card](https://goreportcard.com/badge/github.com/naicoi92/obfuscator)](https://goreportcard.com/report/github.com/naicoi92/obfuscator)
[![GoDoc](https://godoc.org/github.com/naicoi92/obfuscator?status.svg)](https://godoc.org/github.com/naicoi92/obfuscator)

A Go library for JavaScript obfuscation using V8 JavaScript Engine.

[English](README.md) | [Tiếng Việt](README_vi.md)

## Introduction

jsobfuscator-go is a Go library that helps obfuscate JavaScript code by using JavaScript Obfuscator and V8 Engine. This library provides a simple way to protect your JavaScript code from being easily read and understood.

## Installation

```bash
go get github.com/naicoi92/obfuscator
```

## Requirements

- Go 1.24.0 or later
- v8go library (automatically installed when using `go get`)

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"github.com/naicoi92/obfuscator"
)

func main() {
	// Initialize the obfuscator
	obf, err := obfuscator.NewObfuscator()
	if err != nil {
		panic(err)
	}

	// JavaScript code to obfuscate
	jsCode := `function sayHello() { return "Hello World"; }`

	// Perform obfuscation
	obfuscatedCode, err := obf.Obfuscate(jsCode)
	if err != nil {
		panic(err)
	}

	// Print the obfuscated code
	fmt.Println(obfuscatedCode)
}
```

### Important Notes

- JavaScript code cannot contain backtick characters (`) as they are used to wrap the code during processing.
- The obfuscator uses the following default options:
  ```javascript
  const options = {
      compact: (Math.random() < 0.5),
      controlFlowFlattening: true,
      controlFlowFlatteningThreshold: 1,
      numbersToExpressions: true,
      simplify: true,
      stringArrayShuffle: true,
      splitStrings: true,
      stringArrayThreshold: 1
  }
  ```

## Performance Optimization

The library uses a caching mechanism to optimize performance when performing multiple obfuscations. You should reuse the same Obfuscator instance when you need to obfuscate multiple JavaScript code snippets.

```go
obf, _ := obfuscator.NewObfuscator()

// Use the same instance for multiple obfuscations
result1, _ := obf.Obfuscate(jsCode1)
result2, _ := obf.Obfuscate(jsCode2)
result3, _ := obf.Obfuscate(jsCode3)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request or create an Issue on GitHub.

## License

This project is distributed under the MIT license.