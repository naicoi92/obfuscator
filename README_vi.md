# jsobfuscator-go

[![Go Report Card](https://goreportcard.com/badge/github.com/naicoi92/obfuscator)](https://goreportcard.com/report/github.com/naicoi92/obfuscator)
[![GoDoc](https://godoc.org/github.com/naicoi92/obfuscator?status.svg)](https://godoc.org/github.com/naicoi92/obfuscator)

Thư viện Go để làm rối mã JavaScript sử dụng V8 JavaScript Engine.

[English](README.md) | [Tiếng Việt](README_vi.md)

## Giới thiệu

jsobfuscator-go là một thư viện Go giúp làm rối mã JavaScript thông qua việc sử dụng JavaScript Obfuscator và V8 Engine. Thư viện này cung cấp một cách đơn giản để bảo vệ mã JavaScript của bạn khỏi việc bị đọc và hiểu một cách dễ dàng.

## Cài đặt

```bash
go get github.com/naicoi92/obfuscator
```

## Yêu cầu

- Go 1.24.0 trở lên
- Thư viện v8go (được cài đặt tự động khi sử dụng `go get`)

## Cách sử dụng

### Ví dụ cơ bản

```go
package main

import (
	"fmt"
	"github.com/naicoi92/obfuscator"
)

func main() {
	// Khởi tạo obfuscator
	obf, err := obfuscator.NewObfuscator()
	if err != nil {
		panic(err)
	}

	// Mã JavaScript cần làm rối
	jsCode := `function sayHello() { return "Hello World"; }`

	// Thực hiện làm rối mã
	obfuscatedCode, err := obf.Obfuscate(jsCode)
	if err != nil {
		panic(err)
	}

	// In mã đã được làm rối
	fmt.Println(obfuscatedCode)
}
```

### Lưu ý quan trọng

- Mã JavaScript không được chứa ký tự backtick (`) vì nó được sử dụng để bao quanh mã trong quá trình xử lý.
- Obfuscator sử dụng các tùy chọn mặc định sau:
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

## Tối ưu hiệu suất

Thư viện sử dụng cơ chế cache để tối ưu hiệu suất khi thực hiện nhiều lần làm rối mã. Bạn nên tái sử dụng cùng một instance của Obfuscator khi cần làm rối nhiều đoạn mã JavaScript.

```go
obf, _ := obfuscator.NewObfuscator()

// Sử dụng cùng một instance cho nhiều lần làm rối
result1, _ := obf.Obfuscate(jsCode1)
result2, _ := obf.Obfuscate(jsCode2)
result3, _ := obf.Obfuscate(jsCode3)
```

## Đóng góp

Mọi đóng góp đều được hoan nghênh! Vui lòng tạo issue hoặc pull request trên GitHub.

## Giấy phép

Dự án này được phân phối dưới giấy phép MIT.