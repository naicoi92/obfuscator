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

// Default options for JavaScript obfuscation
const defaultOptions = `
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
`

// Obfuscator represents a JavaScript obfuscator instance
type Obfuscator struct {
	CachedData *v8go.CompilerCachedData
}

// NewObfuscator creates and initializes a new JavaScript obfuscator
func NewObfuscator() (*Obfuscator, error) {
	isolate := v8go.NewIsolate()
	defer isolate.Dispose()
	context := v8go.NewContext(isolate)
	defer context.Close()
	o := &Obfuscator{}
	if err := o.setupJSCode(isolate, context, nil); err != nil {
		return nil, fmt.Errorf("failed to setup JS code: %w", err)
	}
	return o, nil
}

// Close releases all resources used by the obfuscator

// setupJSCode loads the JavaScript obfuscator code into the V8 context
func (o *Obfuscator) setupJSCode(
	isolate *v8go.Isolate,
	context *v8go.Context,
	cache *v8go.CompilerCachedData,
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
	if cache != nil {
		opts.CachedData = cache
	}
	script, err := isolate.CompileUnboundScript(code, "obfuscation.js", opts)
	if err != nil {
		return fmt.Errorf("failed to compile script: %w", err)
	}
	if _, err := script.Run(context); err != nil {
		return fmt.Errorf("failed to run script: %w", err)
	}
	if cache == nil {
		o.CachedData = script.CreateCodeCache()
	}
	return nil
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
	if err := o.setupJSCode(isolate, context, o.CachedData); err != nil {
		return "", fmt.Errorf("failed to setup JS code: %w", err)
	}
	codeString := fmt.Sprintf(
		"const code = `%s`; %s ;const obfuscatedCode = JavaScriptObfuscator.obfuscate(code, options).getObfuscatedCode();obfuscatedCode;",
		code,
		defaultOptions,
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
