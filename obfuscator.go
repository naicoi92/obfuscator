package obfuscator

import (
	_ "embed"
	"fmt"

	"rogchap.com/v8go"
)

//go:embed obfuscation.js
var JsCode string

type Obfuscator struct {
	Isolate *v8go.Isolate
	Context *v8go.Context
}

func NewObfuscator() (*Obfuscator, error) {
	o := &Obfuscator{}
	o.Isolate = v8go.NewIsolate()
	o.Context = v8go.NewContext(o.Isolate)
	if err := o.init(); err != nil {
		return nil, err
	}
	return o, nil
}

func (o *Obfuscator) Close() {
	o.Isolate.Dispose()
	o.Context.Close()
}

func (o *Obfuscator) init() error {
	if err := o.setupJSCode(); err != nil {
		return err
	}
	if err := o.checkJSCode(); err != nil {
		return err
	}
	if err := o.setupOptions(); err != nil {
		return err
	}
	return nil
}

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
	if _, err := o.Context.RunScript(code, "obfuscation.js"); err != nil {
		return err
	}
	return nil
}

func (o *Obfuscator) setupOptions() error {
	code := `
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
	if _, err := o.Context.RunScript(code, "options.js"); err != nil {
		return err
	}
	return nil
}

func (o *Obfuscator) Obfuscate(code string) (string, error) {
	codeString := fmt.Sprintf(
		"const code = `%s`;const obfuscatedCode = JavaScriptObfuscator.obfuscate(code, options).getObfuscatedCode();obfuscatedCode;",
		code,
	)
	val, err := o.Context.RunScript(codeString, "obfuscate.js")
	if err != nil {
		return "", err
	}
	obfuscatedCode := val.String()
	if obfuscatedCode == "" {
		return "", fmt.Errorf("obfuscated code is empty")
	}
	return obfuscatedCode, nil
}
