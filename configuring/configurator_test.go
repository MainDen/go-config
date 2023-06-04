package configuring

import (
	"context"
	"fmt"
	"testing"
)

func TestConfigurator_WithLogger(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(args ...interface{}) {
		calls = calls + 1
	}).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls + 1
	}).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(ctx context.Context, args ...interface{}) {
		calls = calls + 1
	}).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(ctx context.Context, format string, args ...interface{}) {
		calls = calls + 1
	}).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(ctx context.Context, args ...interface{}) {
		calls = calls - 1
	}).WithLogger(nil).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	func() {
		defer func() {
			rerr := recover()
			if rerr != nil {
				if err, ok := rerr.(error); !ok || err == nil || err.Error() != "invalid logger function type 'func()' (type should be in ['func(context.Context, string, ...interface {})','func(context.Context, ...interface {})','func(string, ...interface {})','func(...interface {})'])" {
					t.Errorf("expected '%v', was '%v'", "invalid logger function type 'func()' (type should be in ['func(context.Context, string, ...interface {})','func(context.Context, ...interface {})','func(string, ...interface {})','func(...interface {})'])", rerr)
				}
				calls = calls + 1
			}
		}()
		err = NewConfigurator().WithLogger(func() {
			calls = calls - 1
		}).Configure(new(int))
		if err != nil {
			t.Errorf("expected '%v', was '%v'", error(nil), err)
		}
	}()
	if calls != 5 {
		t.Errorf("expected '%v', was '%v'", 5, calls)
	}
}

func TestConfigurator_Secret(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: input: '0' output: '0'" {
			t.Errorf("expected '%v', was '%v'", "configuration: input: '0' output: '0'", message)
		}
		calls = calls + 1
	}).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: input: *secret* output: *secret*" {
			t.Errorf("expected '%v', was '%v'", "configuration: input: *secret* output: *secret*", message)
		}
		calls = calls + 1
	}).Secret().Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	if calls != 2 {
		t.Errorf("expected '%v', was '%v'", 2, calls)
	}
}

func TestConfigurator_WithLogChangesOnly(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(args ...interface{}) {
		calls = calls + 1
	}).WithLogChangesOnly(true).WithDisallowed(0).WithDefault(1).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(args ...interface{}) {
		calls = calls - 1
	}).WithLogChangesOnly(true).WithDisallowed(0).WithCurrent(1).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(args ...interface{}) {
		calls = calls + 1
	}).WithLogChangesOnly(false).WithDisallowed(0).WithCurrent(1).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	if calls != 2 {
		t.Errorf("expected '%v', was '%v'", 2, calls)
	}
}

func TestConfigurator_WithLogValueFormat(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if logMessage := fmt.Sprintf(format, args...); logMessage != "configuration: min: '0' max: '9' allowed: ['0','2','4','6','8'] disallowed: ['1','3','5','7','9'] default: '0' input: '2' output: '2'" {
			t.Errorf("expected '%v', was '%v'", "configuration: min: '0' max: '9' allowed: ['0','2','4','6','8'] disallowed: ['1','3','5','7','9'] default: '0' input: '2' output: '2'", logMessage)
		}
		calls = calls + 1
	}).WithLogValueFormat("'%T (type)'").WithLogValueFormat("").WithMin(0).WithMax(9).WithAllowed(0, 2, 4, 6, 8).WithDisallowed(1, 3, 5, 7, 9).WithDefault(0).WithCurrent(2).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	if calls != 1 {
		t.Errorf("expected '%v', was '%v'", 1, calls)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if logMessage := fmt.Sprintf(format, args...); logMessage != "configuration: min: 'int (type)' max: 'int (type)' allowed: ['int (type)','int (type)','int (type)','int (type)','int (type)'] disallowed: ['int (type)','int (type)','int (type)','int (type)','int (type)'] default: 'int (type)' input: 'int (type)' output: 'int (type)'" {
			t.Errorf("expected '%v', was '%v'", "configuration: min: 'int (type)' max: 'int (type)' allowed: ['int (type)','int (type)','int (type)','int (type)','int (type)'] disallowed: ['int (type)','int (type)','int (type)','int (type)','int (type)'] default: 'int (type)' input: 'int (type)' output: 'int (type)'", logMessage)
		}
		calls = calls + 1
	}).WithLogValueFormat("'%T (type)'").WithMin(0).WithMax(9).WithAllowed(0, 2, 4, 6, 8).WithDisallowed(1, 3, 5, 7, 9).WithDefault(0).WithCurrent(2).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	if calls != 2 {
		t.Errorf("expected '%v', was '%v'", 2, calls)
	}
}

func TestConfigurator_WithContext(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(ctx context.Context, args ...interface{}) {
		if value := ctx.Value(struct{}{}); value != nil {
			t.Errorf("expected '%v', was '%v'", nil, value)
		}
		calls = calls + 1
	}).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(ctx context.Context, args ...interface{}) {
		if value, ok := ctx.Value(struct{}{}).(string); !ok || value != "test context value" {
			t.Errorf("expected '%v', was '%v'", "test context value", value)
		}
		calls = calls + 1
	}).WithContext(context.WithValue(context.TODO(), struct{}{}, "test context value")).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	if calls != 2 {
		t.Errorf("expected '%v', was '%v'", 2, calls)
	}
}

func TestConfigurator_WithName(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: input: '0' output: '0'" {
			t.Errorf("expected '%v', was '%v'", "configuration: input: '0' output: '0'", message)
		}
		calls = calls + 1
	}).WithName("application.value.identifier").WithName("").Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration of 'application.value.identifier': input: '0' output: '0'" {
			t.Errorf("expected '%v', was '%v'", "configuration of 'application.value.identifier': input: '0' output: '0'", message)
		}
		calls = calls + 1
	}).WithName("application.value.identifier").Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithName("application.value.identifier").WithDisallowed(0).Configure(new(int))
	if err == nil || err.Error() != "configuration of 'application.value.identifier' error: target value error: argument should not be in disallowed values ['0']" {
		t.Errorf("expected '%v', was '%v'", "configuration of 'application.value.identifier' error: target value error: argument should not be in disallowed values ['0']", err)
	}
	if calls != 2 {
		t.Errorf("expected '%v', was '%v'", 2, calls)
	}
}

func TestConfigurator_WithMin(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: min: '0' input: '0' output: '0'" {
			t.Errorf("expected '%v', was '%v'", "configuration: min: '0' input: '0' output: '0'", message)
		}
		calls = calls + 1
	}).WithMin(0).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithMin(false).Configure(new(int))
	if err == nil || err.Error() != "configuration error: invalid min value: argument of type 'bool' should be convertible to type 'int'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: invalid min value: argument of type 'bool' should be convertible to type 'int'", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithMin(1).Configure(new(int))
	if err == nil || err.Error() != "configuration error: target value error: argument should be greater than or equal to '1'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value error: argument should be greater than or equal to '1'", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithMin(false).Configure(new(bool))
	if err == nil || err.Error() != "configuration error: target value error: invalid min value: argument of type 'bool' can not be lower than or greater than value of type 'bool'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value error: invalid min value: argument of type 'bool' can not be lower than or greater than value of type 'bool'", err)
	}
	if calls != 1 {
		t.Errorf("expected '%v', was '%v'", 1, calls)
	}
}

func TestConfigurator_WithMax(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: max: '0' input: '0' output: '0'" {
			t.Errorf("expected '%v', was '%v'", "configuration: max: '0' input: '0' output: '0'", message)
		}
		calls = calls + 1
	}).WithMax(0).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithMax(false).Configure(new(int))
	if err == nil || err.Error() != "configuration error: invalid max value: argument of type 'bool' should be convertible to type 'int'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: invalid max value: argument of type 'bool' should be convertible to type 'int'", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithMax(-1).Configure(new(int))
	if err == nil || err.Error() != "configuration error: target value error: argument should be lower than or equal to '-1'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value error: argument should be lower than or equal to '-1'", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithMax(false).Configure(new(bool))
	if err == nil || err.Error() != "configuration error: target value error: invalid max value: argument of type 'bool' can not be lower than or greater than value of type 'bool'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value error: invalid max value: argument of type 'bool' can not be lower than or greater than value of type 'bool'", err)
	}
	if calls != 1 {
		t.Errorf("expected '%v', was '%v'", 1, calls)
	}
}

func TestConfigurator_WithAllowed(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: allowed: ['0','1','2'] input: '0' output: '0'" {
			t.Errorf("expected '%v', was '%v'", "configuration: allowed: ['0','1','2'] input: '0' output: '0'", message)
		}
		calls = calls + 1
	}).WithAllowed(0, 1, 2).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithAllowed(1, 2).Configure(new(int))
	if err == nil || err.Error() != "configuration error: target value error: argument should be in allowed values ['1','2']" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value error: argument should be in allowed values ['1','2']", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithAllowed(0, 1, 2, false).Configure(new(int))
	if err == nil || err.Error() != "configuration error: invalid allowed values: invalid element at index '3': argument of type 'bool' should be convertible to type 'int'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: invalid allowed values: invalid element at index '3': argument of type 'bool' should be convertible to type 'int'", err)
	}
	if calls != 1 {
		t.Errorf("expected '%v', was '%v'", 1, calls)
	}
}

func TestConfigurator_WithDisllowed(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: disallowed: ['1','2'] input: '0' output: '0'" {
			t.Errorf("expected '%v', was '%v'", "configuration: disallowed: ['1','2'] input: '0' output: '0'", message)
		}
		calls = calls + 1
	}).WithDisallowed(1, 2).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithDisallowed(0, 1, 2).Configure(new(int))
	if err == nil || err.Error() != "configuration error: target value error: argument should not be in disallowed values ['0','1','2']" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value error: argument should not be in disallowed values ['0','1','2']", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithDisallowed(1, 2, false).Configure(new(int))
	if err == nil || err.Error() != "configuration error: invalid disallowed values: invalid element at index '2': argument of type 'bool' should be convertible to type 'int'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: disinvalid allowed values: invalid element at index '2': argument of type 'bool' should be convertible to type 'int'", err)
	}
	if calls != 1 {
		t.Errorf("expected '%v', was '%v'", 1, calls)
	}
}

func TestConfigurator_WithValidators(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: input: '[1 2 3]' output: '[1 2 3]'" {
			t.Errorf("expected '%v', was '%v'", "configuration: input: '[1 2 3]' output: '[1 2 3]'", message)
		}
		calls = calls + 1
	}).WithValidators(NewConfigurator().WithAllowed([]int{1, 2, 3})).Configure(&[]int{1, 2, 3})
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithDisallowed(0, 1, 2).Configure(new(int))
	if err == nil || err.Error() != "configuration error: target value error: argument should not be in disallowed values ['0','1','2']" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value error: argument should not be in disallowed values ['0','1','2']", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithDisallowed(1, 2, false).Configure(new(int))
	if err == nil || err.Error() != "configuration error: invalid disallowed values: invalid element at index '2': argument of type 'bool' should be convertible to type 'int'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: disinvalid allowed values: invalid element at index '2': argument of type 'bool' should be convertible to type 'int'", err)
	}
	if calls != 1 {
		t.Errorf("expected '%v', was '%v'", 1, calls)
	}
}

func TestConfigurator_WithDefault(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: allowed: ['1','2'] default: '1' input: '0' output: '1'" {
			t.Errorf("expected '%v', was '%v'", "configuration: allowed: ['1','2'] default: '1' input: '0' output: '1'", message)
		}
		calls = calls + 1
	}).WithAllowed(1, 2).WithDefault(1).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithAllowed(1, 2).WithDefault(false).Configure(new(int))
	if err == nil || err.Error() != "configuration error: invalid default value: argument of type 'bool' should be convertible to type 'int'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: invalid default value: argument of type 'bool' should be convertible to type 'int'", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithAllowed(1, 2).WithDefault(0).Configure(new(int))
	if err == nil || err.Error() != "configuration error: default value error: argument should be in allowed values ['1','2']" {
		t.Errorf("expected '%v', was '%v'", "configuration error: default value error: argument should be in allowed values ['1','2']", err)
	}
	if calls != 1 {
		t.Errorf("expected '%v', was '%v'", 1, calls)
	}
}

func TestConfigurator_WithCurrent(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: allowed: ['1','2'] input: '1' output: '1'" {
			t.Errorf("expected '%v', was '%v'", "configuration: allowed: ['1','2'] input: '1' output: '1'", message)
		}
		calls = calls + 1
	}).WithAllowed(1, 2).WithCurrent(1).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithAllowed(1, 2).WithCurrent(false).Configure(new(int))
	if err == nil || err.Error() != "configuration error: invalid current value: argument of type 'bool' should be convertible to type 'int'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: invalid current value: argument of type 'bool' should be convertible to type 'int'", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithAllowed(1, 2).WithCurrent(0).Configure(new(int))
	if err == nil || err.Error() != "configuration error: target value error: argument should be in allowed values ['1','2']" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value error: argument should be in allowed values ['1','2']", err)
	}
	if calls != 1 {
		t.Errorf("expected '%v', was '%v'", 1, calls)
	}
}

func TestConfigurator_Validate(t *testing.T) {
	err := NewConfigurator().WithAllowed(1, 2).Validate(1)
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithAllowed(1, 2).Validate(nil)
	if err == nil || err.Error() != "validation error: argument should not be nil" {
		t.Errorf("expected '%v', was '%v'", "validation error: argument should not be nil", err)
	}
	err = NewConfigurator().WithAllowed(1, false).Validate(3)
	if err == nil || err.Error() != "validation error: invalid allowed values: invalid element at index '1': argument of type 'bool' should be convertible to type 'int'" {
		t.Errorf("expected '%v', was '%v'", "validation error: invalid allowed values: invalid element at index '1': argument of type 'bool' should be convertible to type 'int'", err)
	}
	err = NewConfigurator().WithAllowed(1, 2).Validate(3)
	if err == nil || err.Error() != "validation error: argument should be in allowed values ['1','2']" {
		t.Errorf("expected '%v', was '%v'", "validation error: argument should be in allowed values ['1','2']", err)
	}
}

func TestConfigurator_Configure(t *testing.T) {
	calls := 0
	err := NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: allowed: ['0','1','2'] input: '0' output: '0'" {
			t.Errorf("expected '%v', was '%v'", "configuration: allowed: ['0','1','2'] input: '0' output: '0'", message)
		}
		calls = calls + 1
	}).WithAllowed(0, 1, 2).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: allowed: ['1','2'] input: '1' output: '1'" {
			t.Errorf("expected '%v', was '%v'", "configuration: allowed: ['1','2'] input: '1' output: '1'", message)
		}
		calls = calls + 1
	}).WithAllowed(1, 2).WithCurrent(1).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).Configure(nil)
	if err == nil || err.Error() != "configuration error: target value is not configurable: argument should be a pointer" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value is not configurable: argument should be a pointer", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).Configure(1)
	if err == nil || err.Error() != "configuration error: target value is not configurable: argument of type 'int' should be a pointer" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value is not configurable: argument of type 'int' should be a pointer", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).Configure((*int)(nil))
	if err == nil || err.Error() != "configuration error: target value is not configurable: argument of type '*int' should not be nil" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value is not configurable: argument of type '*int' should not be nil", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithAllowed(false).Configure(new(int))
	if err == nil || err.Error() != "configuration error: invalid allowed values: invalid element at index '0': argument of type 'bool' should be convertible to type 'int'" {
		t.Errorf("expected '%v', was '%v'", "configuration error: invalid allowed values: invalid element at index '0': argument of type 'bool' should be convertible to type 'int'", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		calls = calls - 1
	}).WithAllowed(1, 2).Configure(new(int))
	if err == nil || err.Error() != "configuration error: target value error: argument should be in allowed values ['1','2']" {
		t.Errorf("expected '%v', was '%v'", "configuration error: target value error: argument should be in allowed values ['1','2']", err)
	}
	err = NewConfigurator().WithLogger(func(format string, args ...interface{}) {
		if message := fmt.Sprintf(format, args...); message != "configuration: allowed: ['1','2'] default: '1' input: '0' output: '1'" {
			t.Errorf("expected '%v', was '%v'", "configuration: allowed: ['1','2'] default: '1' input: '0' output: '1'", message)
		}
		calls = calls + 1
	}).WithLogChangesOnly(true).WithAllowed(1, 2).WithDefault(1).Configure(new(int))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
}
