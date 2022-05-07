package genum

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"unicode"
)

var (
	panicOnEnexistedValues bool
)

func SetPanicOnUnexistedValues(panic bool) {
	panicOnEnexistedValues = panic
}

type BaseType int

var (
	// All writes must be only in initialize state
	// it run in one goroutine
	// then many goroutines can read the map concurrently safe without additional sync
	globalHolders = map[any]any{}
)

type privateType interface {
	~int
}

type EnumValue[T privateType] struct {
	val int
}

func (enum EnumValue[T]) Int() int {
	return enum.val
}

func (enum EnumValue[T]) String() string {
	return enum.string(globalHolders, panicOnEnexistedValues)
}

func (enum EnumValue[T]) string(holders map[any]any, panicOnUnexisted bool) string {
	var zero T

	holder, ok := holders[zero].(*EnumHolder[T])
	if !ok {
		res := fmt.Sprintf("Unexisted holder for type: %v", reflect.TypeOf(zero))
		if panicOnUnexisted {
			panic(res)
		}
		return res
	}

	res, ok := holder.intToString[enum.val]
	if !ok {
		res := fmt.Sprintf("Unexisted string value for: %v of type: %v", enum.val, reflect.TypeOf(zero))
		if panicOnUnexisted {
			panic(res)
		}
		return res
	}
	return res
}

type EnumHolder[T privateType] struct {
	intToString map[int]string
	stringToInt map[string]int
}

// NewHolder must be called ini init state/thread only
func NewHolder[T privateType]() *EnumHolder[T] {
	return newHolder[T](globalHolders)
}

func newHolder[T privateType](m map[any]any) *EnumHolder[T] {
	var zero T
	refType := reflect.TypeOf(zero)
	if refType == reflect.TypeOf(0) {
		panic("deny construct enum direct on int type")
	}

	if !unicode.IsLower([]rune(refType.Name())[0]) {
		panic("privateType must be private (start from lowcase char")
	}

	if _, exist := m[zero]; exist {
		panic(fmt.Sprintf("Holder already exist for type: %v", reflect.TypeOf(zero)))
	}

	holder := &EnumHolder[T]{
		intToString: make(map[int]string),
		stringToInt: make(map[string]int),
	}
	m[zero] = holder
	return holder
}

// New must call only from init state/goroutine
func (h *EnumHolder[T]) New(val int, name string) EnumValue[T] {
	if _, exist := h.intToString[val]; exist {
		var zero T
		panic(fmt.Sprintf("Value already exist: %v for type: %v", val, reflect.TypeOf(zero)))
	}
	if _, exist := h.stringToInt[name]; exist {
		var zero T
		panic(fmt.Sprintf("String name already exist: %q for type: %v", name, reflect.TypeOf(zero)))
	}
	h.intToString[val] = name
	h.stringToInt[name] = val
	return EnumValue[T]{val: val}
}

func (h *EnumHolder[T]) FromInt(val int) (EnumValue[T], error) {
	if _, ok := h.intToString[val]; ok {
		return EnumValue[T]{val: val}, nil
	}
	var zero EnumValue[T]
	return zero, errors.New("int value doesn't exist for the enum")
}

func (h *EnumHolder[T]) FromString(s string) (EnumValue[T], error) {
	if val, ok := h.stringToInt[s]; ok {
		return EnumValue[T]{val: val}, nil
	}
	var zero EnumValue[T]
	return zero, errors.New("string value doesn't exist for the enum")
}

func (h *EnumHolder[T]) ValueToString(v EnumValue[T]) string {
	return h.intToString[v.Int()]
}

// All return all available enum values in int value order
func (h *EnumHolder[T]) All() []EnumValue[T] {
	res := make([]EnumValue[T], 0, len(h.intToString))
	for v := range h.intToString {
		res = append(res, EnumValue[T]{val: v})
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].val < res[j].val
	})
	return res
}
