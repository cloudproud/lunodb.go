package value

import (
	"fmt"
)

func EncodeUUID[T [16]byte](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case [16]byte:
		return append(buf, v[:]...), nil
	default:
		return buf, fmt.Errorf("unsupported uuid type: %T", val)
	}
}
