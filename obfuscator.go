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
	Isolate *v8go.Isolate
	Context *v8go.Context
}

// NewObfuscator creates and initializes a new JavaScript obfuscator
func NewObfuscator() (*Obfuscator, error) {
	isolate := v8go.NewIsolate()
	context := v8go.NewContext(isolate)

	o := &Obfuscator{
		Isolate: isolate,
		Context: context,
	}

	if err := o.init(); err != nil {
		o.Close() // Clean up resources on initialization error
		return nil, fmt.Errorf("failed to initialize obfuscator: %w", err)
	}

	return o, nil
}

// Close releases all resources used by the obfuscator
func (o *Obfuscator) Close() {
	if o.Context != nil {
		o.Context.Close()
	}
	if o.Isolate != nil {
		o.Isolate.Dispose()
	}
}

// init initializes the JavaScript environment for obfuscation
func (o *Obfuscator) init() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"setup JavaScript code", o.setupJSCode},
		{"verify JavaScript Obfuscator", o.checkJSCode},
		{"set up obfuscation options", o.setupOptions},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	return nil
}

// checkJSCode verifies that the JavaScriptObfuscator global is available
func (o *Obfuscator) checkJSCode() error {
	val, err := o.Context.RunScript("typeof JavaScriptObfuscator", "check.js")
	if err != nil {
		return err
	}

	if val.String() != "function" {
		return fmt.Errorf("JavaScriptObfuscator is not defined")
	}

	return nil
}

// setupJSCode loads the JavaScript obfuscator code into the V8 context
func (o *Obfuscator) setupJSCode() error {
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

	_, err := o.Context.RunScript(code, "obfuscation.js")
	return err
}

// setupOptions initializes the default obfuscation options
func (o *Obfuscator) setupOptions() error {
	_, err := o.Context.RunScript(defaultOptions, "options.js")
	return err
}

// Obfuscate transforms the provided JavaScript code using the obfuscator
func (o *Obfuscator) Obfuscate(code string) (string, error) {
	// Escape backticks in the input code to prevent JavaScript template literal issues
	if strings.Contains(code, "`") {
		return "", fmt.Errorf("code cannot contain backtick (`) ")
	}

	codeString := fmt.Sprintf(
		"const code = `%s`;const obfuscatedCode = JavaScriptObfuscator.obfuscate(code, options).getObfuscatedCode();obfuscatedCode;",
		code,
	)

	val, err := o.Context.RunScript(codeString, "obfuscate.js")
	if err != nil {
		return "", fmt.Errorf("obfuscation error: %w", err)
	}

	obfuscatedCode := val.String()
	if obfuscatedCode == "" {
		return "", fmt.Errorf("obfuscated code is empty")
	}

	return obfuscatedCode, nil
}
