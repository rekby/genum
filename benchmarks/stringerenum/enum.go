package stringerenum

//go:generate stringer -type=Enum
type Enum int

const (
	EnumA Enum = iota
	EnumB
	EnumC
)
