package configurator

import (
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func TestBeConfigurable(t *testing.T) {
	err := beConfigurable(nil)
	if err == nil || err.Error() != "argument should be a pointer" {
		t.Errorf("expected '%v', was '%v'", "argument should be a pointer", err)
	}
	var a int
	err = beConfigurable(a)
	if err == nil || err.Error() != "argument of type 'int' should be a pointer" {
		t.Errorf("expected '%v', was '%v'", "argument of type 'int' should be a pointer", err)
	}
	var b *int
	err = beConfigurable(b)
	if err == nil || err.Error() != "argument of type '*int' should not be nil" {
		t.Errorf("expected '%v', was '%v'", "argument of type '*int' should not be nil", err)
	}
	var c *int = new(int)
	err = beConfigurable(c)
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
}

func TestGetValue(t *testing.T) {
	var a int = 1
	err := beConfigurable(&a)
	if err != nil {
		t.Fatalf("expected '%v', was '%v'", error(nil), err)
	}
	value := getValue(&a)
	if value != a {
		t.Errorf("expected '%v', was '%v'", a, value)
	}
}

func TestSetValue(t *testing.T) {
	var a int = 1
	err := beConfigurable(&a)
	if err != nil {
		t.Fatalf("expected '%v', was '%v'", error(nil), err)
	}
	var value int = 2
	setValue(&a, value)
	if a != value {
		t.Errorf("expected '%v', was '%v'", value, a)
	}
}

func TestConvert(t *testing.T) {
	var result interface{}
	var err error
	result, err = convert(int(1), nil)
	if err == nil || err.Error() != "argument should not be nil" {
		t.Errorf("expected '%v', was '%v'", "argument should not be nil", err)
	}
	if result != nil {
		t.Errorf("expected '%v', was '%v'", nil, result)
	}
	result, err = convert(int(1), new(int))
	if err == nil || err.Error() != "argument of type '*int' should be convertible to type 'int'" {
		t.Errorf("expected '%v', was '%v'", "argument of type '*int' should be convertible to type 'int'", err)
	}
	if result != nil {
		t.Errorf("expected '%v', was '%v'", nil, result)
	}
	type t1 int
	type t2 int
	result, err = convert(t1(1), t2(2))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	if resultTyped, ok := result.(t1); !ok || resultTyped != 2 {
		t.Errorf("expected '%v', was '%v'", t1(2), result)
	}
}

func TestConvertNotNil(t *testing.T) {
	var result interface{}
	var err error
	result, err = convertNotNil(int(1), nil)
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	if result != nil {
		t.Errorf("expected '%v', was '%v'", nil, result)
	}
	result, err = convertNotNil(int(1), new(int))
	if err == nil || err.Error() != "argument of type '*int' should be convertible to type 'int'" {
		t.Errorf("expected '%v', was '%v'", "argument of type '*int' should be convertible to type 'int'", err)
	}
	if result != nil {
		t.Errorf("expected '%v', was '%v'", nil, result)
	}
	type t1 int
	type t2 int
	result, err = convertNotNil(t1(1), t2(2))
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	if resultTyped, ok := result.(t1); !ok || resultTyped != 2 {
		t.Errorf("expected '%v', was '%v'", t1(2), result)
	}
}

func TestConvertArray(t *testing.T) {
	p := new(int)
	err := beConfigurable(p)
	if err != nil {
		t.Fatalf("expected '%v', was '%v'", error(nil), err)
	}
	var result []interface{}
	result, err = convertArray(*p, nil)
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	if result != nil {
		t.Errorf("expected '%v', was '%v'", []interface{}(nil), result)
	}
	result, err = convertArray(*p, []interface{}{})
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	if result == nil || len(result) > 0 {
		t.Errorf("expected '%v', was '%v'", []interface{}{}, result)
	}
	type t1 int
	type t2 int
	var a t1 = 1
	var b t2 = 1
	var c *t1
	var d *t2
	result, err = convertArray(p, []interface{}{&a, nil, &b, int(4)})
	if err == nil || err.Error() != "invalid element at index '1': argument should not be nil" {
		t.Errorf("expected '%v', was '%v'", "invalid element at index '1': argument should not be nil", err)
	}
	if result != nil {
		t.Errorf("expected '%v', was '%v'", []interface{}(nil), result)
	}
	result, err = convertArray(p, []interface{}{&a, int(2), &b, nil})
	if err == nil || err.Error() != "invalid element at index '1': argument of type 'int' should be convertible to type '*int'" {
		t.Errorf("expected '%v', was '%v'", "invalid element at index '1': argument of type 'int' should be convertible to type '*int'", err)
	}
	if result != nil {
		t.Errorf("expected '%v', was '%v'", []interface{}(nil), result)
	}
	q := new(t1)
	err = beConfigurable(q)
	if err != nil {
		t.Fatalf("expected '%v', was '%v'", error(nil), err)
	}
	result, err = convertArray(q, []interface{}{&a, &b, c, d})
	if err != nil {
		t.Errorf("expected '%v', was '%v'", error(nil), err)
	}
	expected := []interface{}{&a, (*t1)(unsafe.Pointer(&b)), c, (*t1)(unsafe.Pointer(d))}
	if len(result) != 4 || result[0] != expected[0] || result[1] != expected[1] || result[2] != expected[2] || result[3] != expected[3] {
		t.Errorf("expected '%v', was '%v'", expected, result)
	}
}

type testType struct {
	result bool
}

func (t1 testType) Foo(t2 testType) bool { return t1.result && t2.result }

func (t1 testType) Bar(t2 *testType) bool { return t1.result && t2.result }

func (t1 *testType) Baz(t2 testType) bool { return t1.result && t2.result }

func (t1 *testType) Bat(t2 *testType) bool { return t1.result && t2.result }

func (testType) Quuux() {}

func (*testType) Quuuux() {}

func TestCallMethodBool(t *testing.T) {
	type testCase struct {
		TestCase             string
		ExpectedErrorMessage string
		ExpectedResult       bool
		MethodName           string
		Target               interface{}
		Value                interface{}
	}
	testCases := []testCase{
		{
			TestCase:             "Nil receiver",
			ExpectedErrorMessage: "method '' should not be called with nil receiver",
			Target:               (*testType)(nil),
			Value:                &testType{},
		},
		{
			TestCase:             "Nil value",
			ExpectedErrorMessage: "method '' should not be called with nil argument",
			Target:               &testType{},
			Value:                (*testType)(nil),
		},
		{
			TestCase:             "No method",
			ExpectedErrorMessage: "type 'configurator.testType' has no method 'Invalid'",
			MethodName:           "Invalid",
			Target:               &testType{},
			Value:                &testType{},
		},
		{
			TestCase:             "Invalid method of value",
			ExpectedErrorMessage: "method 'Quuux' of type 'configurator.testType' should accept argument of type 'configurator.testType' or '*configurator.testType' and should return result of type 'bool'",
			MethodName:           "Quuux",
			Target:               &testType{},
			Value:                &testType{},
		},
		{
			TestCase:             "Invalid method of pointer",
			ExpectedErrorMessage: "method 'Quuuux' of type '*configurator.testType' should accept argument of type 'configurator.testType' or '*configurator.testType' and should return result of type 'bool'",
			MethodName:           "Quuuux",
			Target:               &testType{},
			Value:                &testType{},
		},
		{
			TestCase:       "Method with value receiver and value argument",
			ExpectedResult: true,
			MethodName:     "Foo",
			Target:         testType{result: true},
			Value:          testType{result: true},
		},
		{
			TestCase:       "Method with value receiver and pointer argument",
			ExpectedResult: true,
			MethodName:     "Bar",
			Target:         testType{result: true},
			Value:          testType{result: true},
		},
		{
			TestCase:       "Method with pointer receiver and value argument",
			ExpectedResult: true,
			MethodName:     "Baz",
			Target:         testType{result: true},
			Value:          testType{result: true},
		},
		{
			TestCase:       "Method with pointer receiver and pointer argument",
			ExpectedResult: true,
			MethodName:     "Bat",
			Target:         testType{result: true},
			Value:          testType{result: true},
		},
	}
	for _, tc := range testCases {
		rTarget, rValue := indirect(tc.Target, tc.Value)
		if result, err := callMethodBool(tc.MethodName, rTarget, rValue); (err == nil) != (tc.ExpectedErrorMessage == "") || err != nil && err.Error() != tc.ExpectedErrorMessage || result != tc.ExpectedResult {
			t.Errorf("TestCase '%v': expected '%v', was '%v'", tc.TestCase, []interface{}{tc.ExpectedResult, tc.ExpectedErrorMessage}, []interface{}{result, err})
		}
	}
}

func TestIndirect(t *testing.T) {
	i1, i2 := 1, 2
	pi1, pi2 := &i1, &i2
	ppi1, ppi2 := &pi1, &pi2
	if ri1, ri2 := indirect(i1, i2); ri1.Kind() != reflect.Int && ri2.Kind() != reflect.Int {
		t.Errorf("expected '%v', was '%v'", reflect.Int, ri1.Kind())
	}
	if ri1, ri2 := indirect(pi1, pi2); ri1.Kind() != reflect.Int && ri2.Kind() != reflect.Int {
		t.Errorf("expected '%v', was '%v'", reflect.Int, ri1.Kind())
	}
	if ri1, ri2 := indirect(ppi1, ppi2); ri1.Kind() != reflect.Int && ri2.Kind() != reflect.Int {
		t.Errorf("expected '%v', was '%v'", reflect.Int, ri1.Kind())
	}
	if ri1, ri2 := indirect(new(*int), new(*int)); ri1.Kind() != reflect.Ptr && ri2.Kind() != reflect.Ptr {
		t.Errorf("expected '%v', was '%v'", reflect.Ptr, ri1.Kind())
	}
}

func TestEqual(t *testing.T) {
	if result := equal(1, 1); !result {
		t.Errorf("expected '%v', was '%v'", true, result)
	}
	if result := equal(1, 2); result {
		t.Errorf("expected '%v', was '%v'", false, result)
	}
	if result := equal(new(int), new(int)); !result {
		t.Errorf("expected '%v', was '%v'", true, result)
	}
	t1, t2 := time.Unix(0, 1).UTC(), time.Unix(0, 1).In(time.FixedZone("", int(time.Hour)))
	if result := equal(t1, t2); !result {
		t.Errorf("expected '%v', was '%v'", true, result)
	}
	pt1, pt2 := &t1, &t2
	if result := equal(pt1, pt2); !result {
		t.Errorf("expected '%v', was '%v'", true, result)
	}
	if result := equal(&pt1, &pt2); !result {
		t.Errorf("expected '%v', was '%v'", true, result)
	}
	if result := equal(func() {}, func() {}); result {
		t.Errorf("expected '%v', was '%v'", false, result)
	}
}

func TestHasEqual(t *testing.T) {
	if result := hasEqual(2, []interface{}{1, 2, 3}); !result {
		t.Errorf("expected '%v', was '%v'", true, result)
	}
	if result := hasEqual(4, []interface{}{1, 2, 3}); result {
		t.Errorf("expected '%v', was '%v'", false, result)
	}
	if result := hasEqual(1, nil); result {
		t.Errorf("expected '%v', was '%v'", false, result)
	}
	if result := hasEqual(1, []interface{}{}); result {
		t.Errorf("expected '%v', was '%v'", false, result)
	}
}

func TestCompareByMethod(t *testing.T) {
	ri1, ri2 := indirect(1, 1)
	if result, ok := compareByMethod("Equal", ri1, ri2); ok || result != 0 {
		t.Errorf("expected '%v', was '%v'", []interface{}{0, false}, []interface{}{result, ok})
	}
	t1, t2 := time.Unix(0, 1).UTC(), time.Unix(0, 1).In(time.FixedZone("", int(time.Hour)))
	rt1, rt2 := indirect(t1, t2)
	if result, ok := compareByMethod("Before", rt1, rt2); !ok || result != 0 {
		t.Errorf("expected '%v', was '%v'", []interface{}{0, true}, []interface{}{result, ok})
	}
	if result, ok := compareByMethod("After", rt1, rt2); !ok || result != 0 {
		t.Errorf("expected '%v', was '%v'", []interface{}{0, true}, []interface{}{result, ok})
	}
	if result, ok := compareByMethod("Equal", rt1, rt2); ok || result != 0 {
		t.Errorf("expected '%v', was '%v'", []interface{}{0, false}, []interface{}{result, ok})
	}
	if result, ok := compareByMethod("Invalid", rt1, rt2); ok || result != 0 {
		t.Errorf("expected '%v', was '%v'", []interface{}{0, false}, []interface{}{result, ok})
	}
	t1, t2 = time.Unix(0, 1).UTC(), time.Unix(0, 2).In(time.FixedZone("", int(time.Hour)))
	rt1, rt2 = indirect(t1, t2)
	if result, ok := compareByMethod("Before", rt1, rt2); !ok || result != -1 {
		t.Errorf("expected '%v', was '%v'", []interface{}{-1, true}, []interface{}{result, ok})
	}
	if result, ok := compareByMethod("After", rt1, rt2); !ok || result != 1 {
		t.Errorf("expected '%v', was '%v'", []interface{}{1, true}, []interface{}{result, ok})
	}
	if result, ok := compareByMethod("Equal", rt1, rt2); !ok || result != 0 {
		t.Errorf("expected '%v', was '%v'", []interface{}{0, true}, []interface{}{result, ok})
	}
	if result, ok := compareByMethod("Invalid", rt1, rt2); ok || result != 0 {
		t.Errorf("expected '%v', was '%v'", []interface{}{0, false}, []interface{}{result, ok})
	}
}

func TestCompareByMethods(t *testing.T) {
	ri1, ri2 := indirect(1, 1)
	if result, ok := compareByMethods([]string{"Same", "Equal"}, ri1, ri2); ok || result != 0 {
		t.Errorf("expected '%v', was '%v'", []interface{}{0, false}, []interface{}{result, ok})
	}
	t1, t2 := time.Unix(0, 1).UTC(), time.Unix(0, 2).In(time.FixedZone("", int(time.Hour)))
	rt1, rt2 := indirect(t1, t2)
	if result, ok := compareByMethods([]string{"Lower", "Before"}, rt1, rt2); !ok || result != -1 {
		t.Errorf("expected '%v', was '%v'", []interface{}{-1, true}, []interface{}{result, ok})
	}
	if result, ok := compareByMethods([]string{"Greater", "After"}, rt2, rt1); !ok || result != -1 {
		t.Errorf("expected '%v', was '%v'", []interface{}{-1, true}, []interface{}{result, ok})
	}
}

type greaterInt int

func (i1 greaterInt) Greater(i2 greaterInt) bool {
	return i1 > i2
}

type lowerInt int

func (i1 lowerInt) Lower(i2 lowerInt) bool {
	return i1 < i2
}

func TestCompare(t *testing.T) {
	type testCase struct {
		TestCase             string
		ExpectedErrorMessage string
		ExpectedResult       int
		Target               interface{}
		Value                interface{}
	}
	type tfloat32 float32
	type tfloat64 float64
	type tint int
	type tint8 int8
	type tint16 int16
	type tint32 int32
	type tint64 int64
	type tstring string
	type tuint uint
	type tuint8 uint8
	type tuint16 uint16
	type tuint32 uint32
	type tuint64 uint64
	testCases := []testCase{
		{
			TestCase:             "bool: error",
			ExpectedErrorMessage: "argument of type 'bool' can not be lower than or greater than value of type 'bool'",
			Target:               true,
			Value:                false,
		},
		{
			TestCase:       "lowerInt: -1",
			ExpectedResult: -1,
			Target:         lowerInt(1),
			Value:          lowerInt(2),
		},
		{
			TestCase:       "lowerInt: 0",
			ExpectedResult: 0,
			Target:         lowerInt(2),
			Value:          lowerInt(2),
		},
		{
			TestCase:       "lowerInt: 1",
			ExpectedResult: 1,
			Target:         lowerInt(2),
			Value:          lowerInt(1),
		},
		{
			TestCase:       "greaterInt: -1",
			ExpectedResult: -1,
			Target:         greaterInt(1),
			Value:          greaterInt(2),
		},
		{
			TestCase:       "greaterInt: 0",
			ExpectedResult: 0,
			Target:         greaterInt(2),
			Value:          greaterInt(2),
		},
		{
			TestCase:       "greaterInt: 1",
			ExpectedResult: 1,
			Target:         greaterInt(2),
			Value:          greaterInt(1),
		},
		{
			TestCase:       "float32: -1",
			ExpectedResult: -1,
			Target:         tfloat32(1),
			Value:          float32(2),
		},
		{
			TestCase:       "float32: 0",
			ExpectedResult: 0,
			Target:         tfloat32(2),
			Value:          float32(2),
		},
		{
			TestCase:       "float32: 1",
			ExpectedResult: 1,
			Target:         tfloat32(2),
			Value:          float32(1),
		},
		{
			TestCase:       "float64: -1",
			ExpectedResult: -1,
			Target:         tfloat64(1),
			Value:          float64(2),
		},
		{
			TestCase:       "float64: 0",
			ExpectedResult: 0,
			Target:         tfloat64(2),
			Value:          float64(2),
		},
		{
			TestCase:       "float64: 1",
			ExpectedResult: 1,
			Target:         tfloat64(2),
			Value:          float64(1),
		},
		{
			TestCase:       "int: -1",
			ExpectedResult: -1,
			Target:         tint(1),
			Value:          int(2),
		},
		{
			TestCase:       "int: 0",
			ExpectedResult: 0,
			Target:         tint(2),
			Value:          int(2),
		},
		{
			TestCase:       "int: 1",
			ExpectedResult: 1,
			Target:         tint(2),
			Value:          int(1),
		},
		{
			TestCase:       "int8: -1",
			ExpectedResult: -1,
			Target:         tint8(1),
			Value:          int8(2),
		},
		{
			TestCase:       "int8: 0",
			ExpectedResult: 0,
			Target:         tint8(2),
			Value:          int8(2),
		},
		{
			TestCase:       "int8: 1",
			ExpectedResult: 1,
			Target:         tint8(2),
			Value:          int8(1),
		},
		{
			TestCase:       "int16: -1",
			ExpectedResult: -1,
			Target:         tint16(1),
			Value:          int16(2),
		},
		{
			TestCase:       "int16: 0",
			ExpectedResult: 0,
			Target:         tint16(2),
			Value:          int16(2),
		},
		{
			TestCase:       "int16: 1",
			ExpectedResult: 1,
			Target:         tint16(2),
			Value:          int16(1),
		},
		{
			TestCase:       "int32: -1",
			ExpectedResult: -1,
			Target:         tint32(1),
			Value:          int32(2),
		},
		{
			TestCase:       "int32: 0",
			ExpectedResult: 0,
			Target:         tint32(2),
			Value:          int32(2),
		},
		{
			TestCase:       "int32: 1",
			ExpectedResult: 1,
			Target:         tint32(2),
			Value:          int32(1),
		},
		{
			TestCase:       "int64: -1",
			ExpectedResult: -1,
			Target:         tint64(1),
			Value:          int64(2),
		},
		{
			TestCase:       "int64: 0",
			ExpectedResult: 0,
			Target:         tint64(2),
			Value:          int64(2),
		},
		{
			TestCase:       "int64: 1",
			ExpectedResult: 1,
			Target:         tint64(2),
			Value:          int64(1),
		},
		{
			TestCase:       "string: -1",
			ExpectedResult: -1,
			Target:         tstring("1"),
			Value:          string("2"),
		},
		{
			TestCase:       "string: -1",
			ExpectedResult: -1,
			Target:         tstring("1"),
			Value:          string("11"),
		},
		{
			TestCase:       "string: 0",
			ExpectedResult: 0,
			Target:         tstring("2"),
			Value:          string("2"),
		},
		{
			TestCase:       "string: 1",
			ExpectedResult: 1,
			Target:         tstring("2"),
			Value:          string("1"),
		},
		{
			TestCase:       "string: 1",
			ExpectedResult: 1,
			Target:         tstring("11"),
			Value:          string("1"),
		},
		{
			TestCase:       "uint: -1",
			ExpectedResult: -1,
			Target:         tuint(1),
			Value:          uint(2),
		},
		{
			TestCase:       "uint: 0",
			ExpectedResult: 0,
			Target:         tuint(2),
			Value:          uint(2),
		},
		{
			TestCase:       "uint: 1",
			ExpectedResult: 1,
			Target:         tuint(2),
			Value:          uint(1),
		},
		{
			TestCase:       "uint8: -1",
			ExpectedResult: -1,
			Target:         tuint8(1),
			Value:          uint8(2),
		},
		{
			TestCase:       "uint8: 0",
			ExpectedResult: 0,
			Target:         tuint8(2),
			Value:          uint8(2),
		},
		{
			TestCase:       "uint8: 1",
			ExpectedResult: 1,
			Target:         tuint8(2),
			Value:          uint8(1),
		},
		{
			TestCase:       "uint16: -1",
			ExpectedResult: -1,
			Target:         tuint16(1),
			Value:          uint16(2),
		},
		{
			TestCase:       "uint16: 0",
			ExpectedResult: 0,
			Target:         tuint16(2),
			Value:          uint16(2),
		},
		{
			TestCase:       "uint16: 1",
			ExpectedResult: 1,
			Target:         tuint16(2),
			Value:          uint16(1),
		},
		{
			TestCase:       "uint32: -1",
			ExpectedResult: -1,
			Target:         tuint32(1),
			Value:          uint32(2),
		},
		{
			TestCase:       "uint32: 0",
			ExpectedResult: 0,
			Target:         tuint32(2),
			Value:          uint32(2),
		},
		{
			TestCase:       "uint32: 1",
			ExpectedResult: 1,
			Target:         tuint32(2),
			Value:          uint32(1),
		},
		{
			TestCase:       "uint64: -1",
			ExpectedResult: -1,
			Target:         tuint64(1),
			Value:          uint64(2),
		},
		{
			TestCase:       "uint64: 0",
			ExpectedResult: 0,
			Target:         tuint64(2),
			Value:          uint64(2),
		},
		{
			TestCase:       "uint64: 1",
			ExpectedResult: 1,
			Target:         tuint64(2),
			Value:          uint64(1),
		},
	}
	for _, tc := range testCases {
		if result, err := compare(tc.Target, tc.Value); (err == nil) != (tc.ExpectedErrorMessage == "") || err != nil && err.Error() != tc.ExpectedErrorMessage || result != tc.ExpectedResult {
			t.Errorf("TestCase '%v': expected '%v', was '%v'", tc.TestCase, []interface{}{tc.ExpectedResult, tc.ExpectedErrorMessage}, []interface{}{result, err})
		}
	}
}
