package configurator

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	EqualMethodNames   = []string{"Equal"}
	GreaterMethodNames = []string{"Greater", "After"}
	LowerMethodNames   = []string{"Lower", "Before"}
)

func beConfigurable(targetPointer interface{}) error {
	rPointer := reflect.ValueOf(targetPointer)
	if !rPointer.IsValid() {
		return errors.New("argument should be a pointer")
	}
	if rPointer.Kind() != reflect.Ptr {
		return fmt.Errorf("argument of type '%v' should be a pointer", reflect.TypeOf(targetPointer).String())
	}
	if rPointer.IsNil() {
		return fmt.Errorf("argument of type '%v' should not be nil", reflect.TypeOf(targetPointer).String())
	}
	return nil
}

func getValue(pointer interface{}) interface{} {
	return reflect.ValueOf(pointer).Elem().Interface()
}

func setValue(pointer interface{}, value interface{}) {
	reflect.ValueOf(pointer).Elem().Set(reflect.ValueOf(value))
}

func convert(target interface{}, value interface{}) (interface{}, error) {
	rValue := reflect.ValueOf(value)
	if !rValue.IsValid() {
		return nil, errors.New("argument should not be nil")
	}
	rTargetType := reflect.TypeOf(target)
	rValueType := reflect.TypeOf(value)
	if !rValueType.ConvertibleTo(rTargetType) {
		return nil, fmt.Errorf("argument of type '%v' should be convertible to type '%v'", rValueType.String(), rTargetType.String())
	}
	return rValue.Convert(rTargetType).Interface(), nil
}

func convertNotNil(target interface{}, value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}
	return convert(target, value)
}

func convertArray(target interface{}, values []interface{}) ([]interface{}, error) {
	if values == nil {
		return nil, nil
	}
	result := make([]interface{}, len(values))
	var err error
	for i := range values {
		result[i], err = convert(target, values[i])
		if err != nil {
			return nil, fmt.Errorf("invalid element at index '%v': %v", i, err)
		}
	}
	return result, nil
}

func callMethodBool(methodName string, rTarget reflect.Value, rValue reflect.Value) (bool, error) {
	if rTarget.Kind() == reflect.Ptr {
		if rTarget.IsNil() {
			return false, fmt.Errorf("method '%v' should not be called with nil receiver", methodName)
		}
		return false, fmt.Errorf("method '%v' should not be called with nil argument", methodName)
	}
	rTargetType := rTarget.Type()
	rResultType := reflect.TypeOf(bool(true))
	rTargetPointerType := reflect.PtrTo(rTargetType)
	if rMethod, found := rTargetType.MethodByName(methodName); found {
		if rMethod.Type.IsVariadic() || rMethod.Type.NumIn() != 2 || !rMethod.Type.In(1).AssignableTo(rTargetType) && !rMethod.Type.In(1).AssignableTo(rTargetPointerType) || rMethod.Type.NumOut() != 1 || !rMethod.Type.Out(0).AssignableTo(rResultType) {
			return false, fmt.Errorf("method '%v' of type '%v' should accept argument of type '%v' or '%v' and should return result of type '%v'", methodName, rTargetType.String(), rTargetType.String(), rTargetPointerType.String(), rResultType.String())
		}
		if rMethod.Type.In(1).AssignableTo(rTargetType) {
			return rMethod.Func.Call([]reflect.Value{rTarget, rValue})[0].Bool(), nil
		}
		rValuePointer := reflect.New(rTargetType)
		rValuePointer.Elem().Set(rValue)
		return rMethod.Func.Call([]reflect.Value{rTarget, rValuePointer})[0].Bool(), nil
	}
	if rMethod, found := rTargetPointerType.MethodByName(methodName); found {
		if rMethod.Type.IsVariadic() || rMethod.Type.NumIn() != 2 || !rMethod.Type.In(1).AssignableTo(rTargetType) && !rMethod.Type.In(1).AssignableTo(rTargetPointerType) || rMethod.Type.NumOut() != 1 || !rMethod.Type.Out(0).AssignableTo(rResultType) {
			return false, fmt.Errorf("method '%v' of type '%v' should accept argument of type '%v' or '%v' and should return result of type '%v'", methodName, rTargetPointerType.String(), rTargetType.String(), rTargetPointerType.String(), rResultType.String())
		}
		rTargetPointer := reflect.New(rTargetType)
		rTargetPointer.Elem().Set(rTarget)
		if rMethod.Type.In(1).AssignableTo(rTargetType) {
			return rMethod.Func.Call([]reflect.Value{rTargetPointer, rValue})[0].Bool(), nil
		}
		rValuePointer := reflect.New(rTargetType)
		rValuePointer.Elem().Set(rValue)
		return rMethod.Func.Call([]reflect.Value{rTargetPointer, rValuePointer})[0].Bool(), nil
	}
	return false, fmt.Errorf("type '%v' has no method '%v'", rTargetType.String(), methodName)
}

func indirect(target interface{}, value interface{}) (reflect.Value, reflect.Value) {
	rTarget := reflect.ValueOf(target)
	rValue := reflect.ValueOf(value)
	for rTarget.Kind() == reflect.Ptr && !rTarget.IsNil() && !rValue.IsNil() {
		rTarget = rTarget.Elem()
		rValue = rValue.Elem()
	}
	return rTarget, rValue
}

func equal(target interface{}, value interface{}) bool {
	rTarget, rValue := indirect(target, value)
	for _, methodName := range EqualMethodNames {
		if equal, err := callMethodBool(methodName, rTarget, rValue); err == nil {
			return equal
		}
	}
	rTargetType := rTarget.Type()
	if rTargetType.Kind() != reflect.Interface && rTargetType.Comparable() && rTarget.Interface() == rValue.Interface() {
		return true
	}
	return reflect.DeepEqual(target, value)
}

func hasEqual(target interface{}, values []interface{}) bool {
	var equalFound bool
	for _, value := range values {
		equalFound = equalFound || equal(target, value)
	}
	return equalFound
}

func compareByMethod(lowerMethodName string, rTarget reflect.Value, rValue reflect.Value) (int, bool) {
	if lower, err := callMethodBool(lowerMethodName, rTarget, rValue); err == nil {
		if greater, err := callMethodBool(lowerMethodName, rValue, rTarget); err == nil {
			if lower && greater {
				return 0, false
			}
			if lower {
				return -1, true
			}
			if greater {
				return 1, true
			}
			return 0, true
		}
	}
	return 0, false
}

func compareByMethods(lowerMethodNames []string, rTarget reflect.Value, rValue reflect.Value) (int, bool) {
	for _, lowerMethodName := range lowerMethodNames {
		if result, ok := compareByMethod(lowerMethodName, rTarget, rValue); ok {
			return result, ok
		}
	}
	return 0, false
}

func compare(target interface{}, value interface{}) (int, error) {
	rTarget, rValue := indirect(target, value)
	if result, ok := compareByMethods(LowerMethodNames, rTarget, rValue); ok {
		return result, nil
	}
	if result, ok := compareByMethods(GreaterMethodNames, rValue, rTarget); ok {
		return result, nil
	}
	switch rTarget.Kind() {
	case reflect.Float32:
		return compareFloat32(rTarget, rValue), nil
	case reflect.Float64:
		return compareFloat64(rTarget, rValue), nil
	case reflect.Int:
		return compareInt(rTarget, rValue), nil
	case reflect.Int8:
		return compareInt8(rTarget, rValue), nil
	case reflect.Int16:
		return compareInt16(rTarget, rValue), nil
	case reflect.Int32:
		return compareInt32(rTarget, rValue), nil
	case reflect.Int64:
		return compareInt64(rTarget, rValue), nil
	case reflect.String:
		return compareString(rTarget, rValue), nil
	case reflect.Uint:
		return compareUint(rTarget, rValue), nil
	case reflect.Uint8:
		return compareUint8(rTarget, rValue), nil
	case reflect.Uint16:
		return compareUint16(rTarget, rValue), nil
	case reflect.Uint32:
		return compareUint32(rTarget, rValue), nil
	case reflect.Uint64:
		return compareUint64(rTarget, rValue), nil
	}
	return 0, fmt.Errorf("argument of type '%v' can not be lower than or greater than value of type '%v'", reflect.TypeOf(target).String(), reflect.TypeOf(value).String())
}

func compareFloat32(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toFloat32Value(rValue1)
	value2 := toFloat32Value(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareFloat64(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toFloat64Value(rValue1)
	value2 := toFloat64Value(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareInt(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toIntValue(rValue1)
	value2 := toIntValue(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareInt8(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toInt8Value(rValue1)
	value2 := toInt8Value(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareInt16(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toInt16Value(rValue1)
	value2 := toInt16Value(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareInt32(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toInt32Value(rValue1)
	value2 := toInt32Value(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareInt64(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toInt64Value(rValue1)
	value2 := toInt64Value(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareString(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toStringValue(rValue1)
	value2 := toStringValue(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareUint(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toUintValue(rValue1)
	value2 := toUintValue(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareUint8(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toUint8Value(rValue1)
	value2 := toUint8Value(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareUint16(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toUint16Value(rValue1)
	value2 := toUint16Value(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareUint32(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toUint32Value(rValue1)
	value2 := toUint32Value(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func compareUint64(rValue1 reflect.Value, rValue2 reflect.Value) int {
	value1 := toUint64Value(rValue1)
	value2 := toUint64Value(rValue2)
	if value1 > value2 {
		return 1
	}
	if value1 < value2 {
		return -1
	}
	return 0
}

func toFloat32Value(rValue reflect.Value) float32 {
	value, ok := rValue.Interface().(float32)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(float32)
	}
	return value
}

func toFloat64Value(rValue reflect.Value) float64 {
	value, ok := rValue.Interface().(float64)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(float64)
	}
	return value
}

func toIntValue(rValue reflect.Value) int {
	value, ok := rValue.Interface().(int)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(int)
	}
	return value
}

func toInt8Value(rValue reflect.Value) int8 {
	value, ok := rValue.Interface().(int8)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(int8)
	}
	return value
}

func toInt16Value(rValue reflect.Value) int16 {
	value, ok := rValue.Interface().(int16)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(int16)
	}
	return value
}

func toInt32Value(rValue reflect.Value) int32 {
	value, ok := rValue.Interface().(int32)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(int32)
	}
	return value
}

func toInt64Value(rValue reflect.Value) int64 {
	value, ok := rValue.Interface().(int64)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(int64)
	}
	return value
}

func toStringValue(rValue reflect.Value) string {
	value, ok := rValue.Interface().(string)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(string)
	}
	return value
}

func toUintValue(rValue reflect.Value) uint {
	value, ok := rValue.Interface().(uint)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(uint)
	}
	return value
}

func toUint8Value(rValue reflect.Value) uint8 {
	value, ok := rValue.Interface().(uint8)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(uint8)
	}
	return value
}

func toUint16Value(rValue reflect.Value) uint16 {
	value, ok := rValue.Interface().(uint16)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(uint16)
	}
	return value
}

func toUint32Value(rValue reflect.Value) uint32 {
	value, ok := rValue.Interface().(uint32)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(uint32)
	}
	return value
}

func toUint64Value(rValue reflect.Value) uint64 {
	value, ok := rValue.Interface().(uint64)
	if !ok {
		value = rValue.Convert(reflect.TypeOf(value)).Interface().(uint64)
	}
	return value
}
