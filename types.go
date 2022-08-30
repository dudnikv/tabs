//

package tabs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	//	"github.com/dudnikv/rf"
)

type Value = uint64

type Table struct {
	cols []Column
	data [][]Value
}

type Column struct {
	cname string
	ctype ValueType
}

type Idented interface {
	Ident() string
}

type IdentedObject struct {
	name string
}

func (p *IdentedObject) Ident() string { return p.name }

var registeredValueTypes = make(map[string]ValueType)

func registerValueType(t ValueType) {
	key := t.Ident()
	if _, ok := registeredValueTypes[key]; ok {
		panic(fmt.Errorf("Value type alredy '%s' exists", key))
	}
	registeredValueTypes[key] = t
}

type ValueType interface {
	Idented
	Parse(string) (Value, error)
	Validate(*string) error
	Emit(Value) string
	Valid(Value) bool
	//	Empty() Value
}

type DefaultValueType struct {
	IdentedObject
	s string
	v Value
}

func (t *DefaultValueType) Emit(v Value) string           { return t.s }
func (t *DefaultValueType) Valid(v Value) bool            { return v == t.v }
func (t *DefaultValueType) Validate(s *string) error      { return nil }
func (t *DefaultValueType) Parse(s string) (Value, error) { return t.v, nil }

var notEmpty error = errors.New("Not empty")

type EmptyValueType struct{}

//func (t *EmptyValueType) Empty() Value        { return 0 }
//func (t *EmptyValueType) HasVoid() bool       { return false }
//func (t *EmptyValueType) IsCompleet() bool    { return true }
//func (t *EmptyValueType) Compleet()           {}
func (t *EmptyValueType) Emit(v Value) string { return "" }
func (t *EmptyValueType) Valid(v Value) bool  { return v == 0 }

func (t *EmptyValueType) Validate(s *string) error {
	if tmp := strings.TrimSpace(*s); tmp == "" {
		*s = tmp
		return nil
	}
	return notEmpty
}

func (t *EmptyValueType) Parse(s string) (Value, error) {
	err := t.Validate(&s)
	return 0, err
}

type TimeValueType struct {
	IdentedObject
	format string
}

func (t *TimeValueType) Emit(v Value) string {
	tt := time.Unix(int64(v), 0)
	return tt.Format(t.format)
}

//func

type EnumValueType struct {
	IdentedObject
	void bool
	cplt bool
	enum map[Value]string
	list map[string]Value
}

func GetValueType(name string) ValueType {
	return registeredValueTypes[name]
}

func NewEnumValueType(name string, list ...string) *EnumValueType {
	var t EnumValueType
	t.name = name
	t.void = true
	t.cplt = false
	t.enum = make(map[Value]string)
	t.list = make(map[string]Value)
	for _, item := range list {
		t.Parse(item)
		//		fmt.Println(t.Parse(item))
	}
	registerValueType(&t)
	return &t
}

func (t *EnumValueType) Empty() Value             { return 0 }
func (t *EnumValueType) HasVoid() bool            { return t.void }
func (t *EnumValueType) Validate(s *string) error { return nil }
func (t *EnumValueType) IsCompleet() bool         { return t.cplt }
func (t *EnumValueType) Compleet()                { t.cplt = true }
func (t *EnumValueType) Emit(v Value) string      { return t.enum[v] }
func (t *EnumValueType) Valid(v Value) bool {
	return v == 0 && t.HasVoid() || v > 0 && v <= uint64(len(t.enum))
}

func (t *EnumValueType) Parse(s string) (Value, error) {
	if s == "" {
		if t.HasVoid() {
			return 0, nil
		}
		return 0, fmt.Errorf("Can't use empty value for %T %s", *t, t.name)
	}
	v, ok := t.list[s]
	if ok {
		return v, nil
	}
	if t.IsCompleet() {
		return 0, fmt.Errorf("Type %s: key '%s' not found", t.name, s)
	}
	if err := t.Validate(&s); err != nil {
		return 0, err
	}
	v = uint64(len(t.list)) + 1
	t.list[s] = v
	t.enum[v] = s
	return v, nil
}

type DataType interface {
	Parse(string) (interface{}, error)
	Emit(interface{}) string
	Empty() interface{}
}

type easVersion = uint64

type easVersionType struct {
	zero easVersion
}

func (t easVersionType) Empty() interface{} { return t.zero }

func (t easVersionType) Valid(v interface{}) bool {
	if x, ok := v.(uint64); ok {
		if x >= 100000000 && x <= 999999999 || x == t.zero {
			return true
		}
	}
	return false
}

func (t easVersionType) Parse(s string) (interface{}, error) {
	if s == "" {
		return t.zero, nil
	}
	s = strings.TrimSpace(s)
	v, err := strconv.ParseUint(strings.Replace(s, ".", "", -1), 10, 64)
	if err != nil {
		goto retErr
	}
	if !t.Valid(v) {
		goto retErr
	}
	return v, nil
retErr:
	return t.zero, fmt.Errorf("Bad EAS version %s", s)
}

func (t easVersionType) Emit(v interface{}) string {
	if !t.Valid(v) {
		panic(fmt.Errorf("Bad EAS version value %v", v))
	}
	//	fmt.Println(v)
	x := v.(uint64)
	if x == t.zero {
		return ""
	}
	s := strconv.FormatUint(v.(uint64), 10)
	return s[:2] + "." + s[2:3] + "." + s[3:4] + "." + s[4:]
}
