package genum

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// MarshalText implements https://pkg.go.dev/encoding#TextMarshaler
func (enum EnumValue[T]) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}

// UnmarshalText implements https://pkg.go.dev/encoding#TextUnmarshaler
func (enum *EnumValue[T]) UnmarshalText(text []byte) error {
	holder, ok := getHolder[T]()
	if !ok {
		return fmt.Errorf("genum: can't get holder for the type %T", enum)
	}
	val, err := holder.FromString(string(text))
	if err != nil {
		return fmt.Errorf("genum: can't unmarshal enum value %q: %w", string(text), err)
	}
	*enum = val
	return nil
}

// MarshalBinary implements https://pkg.go.dev/encoding#BinaryMarshaler
func (enum EnumValue[T]) MarshalBinary() (data []byte, err error) {
	var marshalledVal = int64(enum.val)

	const maxBytesForVarint = 10
	bufArr := make([]byte, maxBytesForVarint)
	varLen := binary.PutVarint(bufArr[:], marshalledVal)
	return bufArr[:varLen], nil
}

func (enum *EnumValue[T]) UnmarshalBinary(data []byte) error {
	val, readBytes := binary.Varint(data)
	if readBytes <= 0 || readBytes != len(data) {
		return errors.New("genum: bad binary format")
	}
	holder, ok := getHolder[T]()
	if !ok {
		return fmt.Errorf("genum: unexpected value type (holder not found): %T", enum)
	}
	newEnum, err := holder.FromInt(int(val))
	if err != nil {
		return err
	}
	*enum = newEnum
	return nil
}

func newVarintBuf() []byte {
	return make([]byte, 8)
}
