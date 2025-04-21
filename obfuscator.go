// Package obfuscator provides functionality to obfuscate JavaScript code
// using JavaScript Obfuscator through v8go.
package obfuscator

import (
	_ "embed"
	"fmt"
	"strings"

	"rogchap.com/v8go"
)

// JsCode contains the embedded JavaScript obfuscator code
//
//go:embed obfuscation.js
var JsCode string

// Obfuscator represents a JavaScript obfuscator instance
type Obfuscator struct {
	CachedData *v8go.CompilerCachedData
	Level      ObfuscationLevel
}

// NewObfuscator creates and initializes a new JavaScript obfuscator
func NewObfuscator() (*Obfuscator, error) {
	isolate := v8go.NewIsolate()
	defer isolate.Dispose()
	context := v8go.NewContext(isolate)
	defer context.Close()
	o := &Obfuscator{
		Level: DefaultLevel,
	}
	if err := o.setupJSCode(isolate, context); err != nil {
		return nil, fmt.Errorf("failed to setup JS code: %w", err)
	}
	return o, nil
}

// Close releases all resources used by the obfuscator

// setupJSCode loads the JavaScript obfuscator code into the V8 context
func (o *Obfuscator) setupJSCode(
	isolate *v8go.Isolate,
	context *v8go.Context,
) error {
	code := fmt.Sprintf(`
  (function() {
    var self = this;
    var window = this;
    var module = {};
    var exports = {};
    module.exports = exports;
    %s
    globalThis.JavaScriptObfuscator = module.exports;
	})()
  `, JsCode)
	opts := v8go.CompileOptions{}
	if o.CachedData != nil {
		opts.CachedData = o.CachedData
	}
	script, err := isolate.CompileUnboundScript(code, "obfuscation.js", opts)
	if err != nil {
		return fmt.Errorf("failed to compile script: %w", err)
	}
	if _, err := script.Run(context); err != nil {
		return fmt.Errorf("failed to run script: %w", err)
	}
	if o.CachedData == nil {
		o.CachedData = script.CreateCodeCache()
	}
	return nil
}

func (o *Obfuscator) SetLevel(level string) {
	o.Level = ObfuscationLevel(level)
}

// Obfuscate transforms the provided JavaScript code using the obfuscator
func (o *Obfuscator) Obfuscate(code string) (string, error) {
	// Escape backticks in the input code to prevent JavaScript template literal issues
	if strings.Contains(code, "`") {
		return "", fmt.Errorf("code cannot contain backtick (`) ")
	}
	isolate := v8go.NewIsolate()
	defer isolate.Dispose()
	context := v8go.NewContext(isolate)
	defer context.Close()
	if err := o.setupJSCode(isolate, context); err != nil {
		return "", fmt.Errorf("failed to setup JS code: %w", err)
	}
	options := getOptions(o.Level)
	codeString := fmt.Sprintf(
		"const code = `%s`; %s ;const obfuscatedCode = JavaScriptObfuscator.obfuscate(code, options).getObfuscatedCode();obfuscatedCode;",
		code,
		options,
	)
	val, err := context.RunScript(codeString, "run.js")
	if err != nil {
		return "", fmt.Errorf("obfuscation error: %w", err)
	}
	obfuscatedCode := val.String()
	if obfuscatedCode == "" {
		return "", fmt.Errorf("obfuscated code is empty")
	}
	return obfuscatedCode, nil
}
