package value

import (
	"encoding/binary"
	"fmt"
)

func EncodeString[T string | *string](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case string:
		buf = binary.BigEndian.AppendUint64(buf, uint64(len(v)))
		return append(buf, []byte(v)...), nil
	case *string:
		if v == nil {
			return buf, nil
		}

		buf = binary.BigEndian.AppendUint64(buf, uint64(len(*v)))
		return append(buf, []byte(*v)...), nil
	default:
		return buf, fmt.Errorf("unsupported string type: %T", val)
	}
}
