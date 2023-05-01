package configurator

import (
	"context"
	"fmt"
	"log"
	"strings"
)

var (
	Default = NewConfigurator().WithLogger(log.Printf)
	Secret  = NewConfigurator().WithLogger(log.Printf).WithLogSecret()
)

type Configurator struct {
	ctx            context.Context
	name           string
	logFn          interface{}
	logChangesOnly bool
	logValueFormat string

	minValue         interface{}
	maxValue         interface{}
	allowedValues    []interface{}
	disallowedValues []interface{}
	defaultValue     interface{}
	currentValue     interface{}
	isSecret         bool
}

func NewConfigurator() Configurator {
	return Configurator{}
}

// WithLogger provides logger function.
// Allowed logFn types:
// func(ctx context.Context, format string, args ...interface{})
// func(ctx context.Context, args ...interface{})
// func(format string, args ...interface{})
// func(args ...interface{})
func (c Configurator) WithLogger(logFn interface{}) Configurator {
	switch logFn := logFn.(type) {
	case func(context.Context, string, ...interface{}):
	case func(context.Context, ...interface{}):
	case func(string, ...interface{}):
	case func(...interface{}):
	case nil:
	default:
		panic(fmt.Errorf("invalid logger function type '%T' (type should be in ['func(context.Context, string, ...interface {})','func(context.Context, ...interface {})','func(string, ...interface {})','func(...interface {})'])", logFn))
	}
	c.logFn = logFn
	return c
}

// WithLogSecret hides input and output values on log message.
func (c Configurator) WithLogSecret() Configurator {
	c.isSecret = true
	return c
}

// WithLogChangesOnly defines when configurator should print log message.
// If set to log changes only configurator prints log message when input value changes.
// Value can be changed with default if invalid.
func (c Configurator) WithLogChangesOnly(logChangesOnly bool) Configurator {
	c.logChangesOnly = logChangesOnly
	return c
}

// WithLogValueFormat provides log format for values on log message. Default value is "'%v'".
// Examples:
// configurator.NewConfigurator().WithLogValueFormat("%v")
// configurator.NewConfigurator().WithLogValueFormat("'%T (type)'")
func (c Configurator) WithLogValueFormat(logValueFormat string) Configurator {
	c.logValueFormat = logValueFormat
	return c
}

func (c Configurator) WithContext(ctx context.Context) Configurator {
	c.ctx = ctx
	return c
}

// WithName provides name of configuration. Helps to find issues in logs.
// It may be common configuration name or specific for configurable value.
func (c Configurator) WithName(name string) Configurator {
	c.name = name
	return c
}

func (c Configurator) WithMin(minValue interface{}) Configurator {
	c.minValue = minValue
	return c
}

func (c Configurator) WithMax(maxValue interface{}) Configurator {
	c.maxValue = maxValue
	return c
}

func (c Configurator) WithAllowed(values ...interface{}) Configurator {
	c.allowedValues = values
	return c
}

func (c Configurator) WithDisallowed(values ...interface{}) Configurator {
	c.disallowedValues = values
	return c
}

func (c Configurator) WithDefault(defaultValue interface{}) Configurator {
	c.defaultValue = defaultValue
	return c
}

func (c Configurator) WithCurrent(currentValue interface{}) Configurator {
	c.currentValue = currentValue
	return c
}

func (c Configurator) log(inputValue interface{}, outputValue interface{}) {
	if c.logFn == nil {
		return
	}
	logValueFormat := "'%v'"
	if c.logValueFormat != "" {
		logValueFormat = c.logValueFormat
	}
	var builder strings.Builder
	var args []interface{}
	_, _ = builder.WriteString("configuration")
	if len(c.name) != 0 {
		_, _ = builder.WriteString(" of '%v'")
		args = append(args, c.name)
	}
	_, _ = builder.WriteString(":")
	if c.minValue != nil {
		_, _ = builder.WriteString(" min: ")
		_, _ = builder.WriteString(logValueFormat)
		args = append(args, c.minValue)
	}
	if c.maxValue != nil {
		_, _ = builder.WriteString(" max: ")
		_, _ = builder.WriteString(logValueFormat)
		args = append(args, c.maxValue)
	}
	if len(c.allowedValues) > 0 {
		_, _ = builder.WriteString(fmt.Sprintf(" allowed: [%v%v]", logValueFormat, strings.Repeat(","+logValueFormat, len(c.allowedValues)-1)))
		args = append(args, c.allowedValues...)
	}
	if len(c.disallowedValues) > 0 {
		_, _ = builder.WriteString(fmt.Sprintf(" disallowed: [%v%v]", logValueFormat, strings.Repeat(","+logValueFormat, len(c.disallowedValues)-1)))
		args = append(args, c.disallowedValues...)
	}
	if c.defaultValue != nil {
		_, _ = builder.WriteString(" default: ")
		_, _ = builder.WriteString(logValueFormat)
		args = append(args, c.defaultValue)
	}
	if c.isSecret {
		logValueFormat = "*secret*"
	}
	_, _ = builder.WriteString(" input: ")
	_, _ = builder.WriteString(logValueFormat)
	_, _ = builder.WriteString(" output: ")
	_, _ = builder.WriteString(logValueFormat)
	if !c.isSecret {
		args = append(args, inputValue, outputValue)
	}
	ctx := c.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	switch logFn := c.logFn.(type) {
	case func(context.Context, string, ...interface{}):
		logFn(ctx, builder.String(), args...)
	case func(context.Context, ...interface{}):
		logFn(ctx, fmt.Sprintf(builder.String(), args...))
	case func(string, ...interface{}):
		logFn(builder.String(), args...)
	case func(...interface{}):
		logFn(fmt.Sprintf(builder.String(), args...))
	}
}

func (c Configurator) wrapError(actionName string, err error) error {
	if err == nil {
		return nil
	}
	if c.name == "" {
		return fmt.Errorf("%v error: %v", actionName, err)
	}
	return fmt.Errorf("%v of '%v' error: %v", actionName, c.name, err)
}

func (c Configurator) convert(target interface{}) (Configurator, error) {
	var err error
	if c.currentValue, err = convertNotNil(target, c.currentValue); err != nil {
		return c, fmt.Errorf("invalid current value: %v", err)
	}
	if c.defaultValue, err = convertNotNil(target, c.defaultValue); err != nil {
		return c, fmt.Errorf("invalid default value: %v", err)
	}
	if c.minValue, err = convertNotNil(target, c.minValue); err != nil {
		return c, fmt.Errorf("invalid min value: %v", err)
	}
	if c.maxValue, err = convertNotNil(target, c.maxValue); err != nil {
		return c, fmt.Errorf("invalid max value: %v", err)
	}
	if c.allowedValues, err = convertArray(target, c.allowedValues); err != nil {
		return c, fmt.Errorf("invalid allowed values: %v", err)
	}
	if c.disallowedValues, err = convertArray(target, c.disallowedValues); err != nil {
		return c, fmt.Errorf("invalid disallowed values: %v", err)
	}
	return c, nil
}

func (c Configurator) validate(target interface{}) error {
	if c.minValue != nil {
		comparisonResult, err := compare(target, c.minValue)
		if err != nil {
			return fmt.Errorf("invalid min value: %v", err)
		}
		if comparisonResult == -1 {
			return fmt.Errorf("argument should be greater than or equal to '%v'", c.minValue)
		}
	}
	if c.maxValue != nil {
		comparisonResult, err := compare(target, c.maxValue)
		if err != nil {
			return fmt.Errorf("invalid max value: %v", err)
		}
		if comparisonResult == 1 {
			return fmt.Errorf("argument should be lower than or equal to '%v'", c.maxValue)
		}
	}
	if len(c.allowedValues) > 0 {
		if !hasEqual(target, c.allowedValues) {
			return fmt.Errorf(fmt.Sprintf("argument should be in allowed values [%v%v]", "'%v'", strings.Repeat(",'%v'", len(c.allowedValues)-1)), c.allowedValues...)
		}
	}
	if len(c.disallowedValues) > 0 {
		if hasEqual(target, c.disallowedValues) {
			return fmt.Errorf(fmt.Sprintf("argument should not be in disallowed values [%v%v]", "'%v'", strings.Repeat(",'%v'", len(c.disallowedValues)-1)), c.disallowedValues...)
		}
	}
	return nil
}

func (c Configurator) Validate(target interface{}) error {
	var err error
	if target, err = convert(target, target); err != nil {
		return c.wrapError("validation", err)
	}
	if c, err = c.convert(target); err != nil {
		return c.wrapError("validation", err)
	}
	return c.wrapError("validation", c.validate(target))
}

func (c Configurator) configure(target interface{}) (interface{}, error) {
	var err error
	if c.defaultValue != nil {
		if err = c.validate(c.defaultValue); err != nil {
			return nil, fmt.Errorf("default value error: %v", err)
		}
	}
	if err = c.validate(target); err != nil {
		if c.defaultValue != nil {
			return c.defaultValue, nil
		}
		return nil, fmt.Errorf("target value error: %v", err)
	}
	return target, nil
}

func (c Configurator) Configure(targetPointer interface{}) error {
	var err error
	if err = beConfigurable(targetPointer); err != nil {
		return c.wrapError("configuration", fmt.Errorf("target value is not configurable: %v", err))
	}
	target := getValue(targetPointer)
	if c, err = c.convert(target); err != nil {
		return c.wrapError("configuration", err)
	}
	if c.currentValue != nil {
		target = c.currentValue
	}
	result, err := c.configure(target)
	if err != nil {
		return c.wrapError("configuration", err)
	}
	if !c.logChangesOnly || !equal(target, result) {
		c.log(target, result)
	}
	setValue(targetPointer, result)
	return nil
}
